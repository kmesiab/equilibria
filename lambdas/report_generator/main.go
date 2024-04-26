package report_generator

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/utils"
	"github.com/kmesiab/equilibria/lambdas/models"
)

const (
	MaxNewMemories           = 1000
	MaxOldMemories           = 1000
	MinimumMemoriesForReport = 50
)

func (p *PatientReportGeneratorLambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "POST":

		return p.Create(request)

		// Enable cors Preflight
	case "OPTIONS":
		return events.APIGatewayProxyResponse{
			Headers:    config.DefaultHttpHeaders,
			StatusCode: http.StatusOK,
		}, nil
	default:

		return lib.RespondWithError("Unsupported HTTP method", nil, http.StatusMethodNotAllowed)
	}
}

func (p *PatientReportGeneratorLambdaHandler) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		repRequest PatientReportGenerationRequest
		patient    *models.User
		memories   *[]models.Message
		report     string
		err        error
	)

	if err = json.Unmarshal([]byte(request.Body), &repRequest); err != nil {

		return lib.RespondWithError("Invalid request JSON.", err, http.StatusBadRequest)
	}

	// Request for a report from system user is ignored.
	if repRequest.PatientID == models.GetSystemUser().ID {

		return lib.RespondWithError("Invalid request.", err, http.StatusBadRequest)
	}

	// Get the patient
	if patient, err = p.UserService.GetUserByID(repRequest.PatientID); err != nil {

		return lib.RespondWithError("Unknown patient ID.", err, http.StatusBadRequest)
	}

	// Generate the prompt for the report
	prompt := GeneratePrompt(repRequest.ReportType, *patient)

	if memories, err = p.MemoryService.GetMemories(patient); err != nil {

		return lib.RespondWithError("Could not retrieve memories.", err, http.StatusInternalServerError)
	}

	// Retrieve all patient memories
	if len(*memories) < MinimumMemoriesForReport {

		msg := log.New("Not enough messages to produce a report.  %d of %d available.",
			len(*memories), MinimumMemoriesForReport).Write()

		return lib.RespondWithError(msg, nil, http.StatusInternalServerError)
	}

	// Generate the report
	if report, err = p.CompletionService.GetCompletion("", prompt, memories); err != nil {

		return lib.RespondWithError("Could not generate a completion.", err, http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Headers:    config.DefaultHttpHeaders,
		StatusCode: http.StatusCreated,
		Body:       report,
	}, nil
}

func Unused() {
	log.New("Report Generator Lambda booting.....").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config")
	}

	database := db.Get(cfg)

	if err := utils.PingDatabase(database); err != nil {
		log.New("Error pinging database").AddError(err).Log()

		return
	}

	if err := utils.PingGoogle(); err != nil {
		log.New("Error pinging Google. Possible bad internet connection.").
			AddError(err).Log()

		return
	}

	memSvc := message.NewMemoryService(
		message.NewMessageRepository(database), MaxNewMemories+MaxOldMemories,
	)

	llmSvc := &ai.OpenAICompletionService{
		RemoveEmojis: false,
	}

	handler := &PatientReportGeneratorLambdaHandler{
		MemoryService:     memSvc,
		CompletionService: llmSvc,
	}

	handler.Init(database)

	log.New("Report Generator Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
