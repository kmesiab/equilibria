package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/emotions"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/nrclex"
	"github.com/kmesiab/equilibria/lambdas/lib/sqs"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/lib/utils"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// How many memories to include in the prompt
// Turns into 75 random user messages?
const maxMemories = 1

// How many immediately previous messages to include in the prompt
const maxLastFewMessages = 250

// How many memories you have to have before we consider you an 'existing'
// user, so the model treats you like it knows you well.
const newUserMemoryCount = 5

type SendSMSLambdaHandler struct {
	lib.LambdaHandler

	MaxMemories        int
	MaxLastFewMemories int

	MemoryService     *message.MemoryService
	CompletionService ai.CompletionServiceInterface
	NRCLexService     *emotions.NRCLexService
}

func (h *SendSMSLambdaHandler) HandleRequest(sqsEvent events.SQSEvent) {

	defer func() {
		if r := recover(); r != nil {
			log.New("Panic while processing event: %v\nStack trace:\n%s", r, debug.Stack()).Log()
		}
	}()

	var (
		err      error
		nowInUTC = time.Now().UTC()

		msg         models.Message
		eventRecord sqs.SQSEventRecord
	)

	if err = ValidateEvent(&sqsEvent); err != nil {
		log.New("Error validating the event").AddError(err).Log()

		return
	}

	if len(sqsEvent.Records) == 0 {

		log.New("No records found in the event. Shutting down.").AddError(err).Log()

		return
	}

	// Unpack the SNS Event Record
	var event = sqsEvent.Records[0]
	if err = json.Unmarshal([]byte(event.Body), &eventRecord); err != nil {
		log.New("Error unmarshalling event record").
			AddError(err).Log()

		return
	}

	// Unpack the message from the event record
	if err = json.Unmarshal([]byte(eventRecord.Message), &msg); err != nil {
		log.New("Error unmarshalling message from event record").
			AddError(err).Log()

		return
	}

	// Get the sender's user account
	recipient, err := h.UserService.GetUserByID(msg.FromUserID)

	if err != nil {

		theMessage, _ := json.Marshal(msg)

		log.New("Unmarshalled Message from Body: %s", theMessage).Log()

		log.New("Could not locate user %d.  Rejecting.", msg.FromUserID).
			AddSQSEvent(&event).AddError(err).Log()

		return
	}

	// Make sure the sender is not the system user
	if recipient.ID == models.GetSystemUser().ID {
		log.New("Message is from system user. Aborting.").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	log.New("Starting response for %s", recipient.PhoneNumber).
		AddUser(recipient).AddSQSEvent(&event).AddMessage(&msg).Log()

	// Get the memories for the user
	memories, err := h.GetMemories(recipient, event, msg)

	log.New("Attaching %d memories", len(memories)).
		AddUser(recipient).
		Log()

	var promptModifier = ExistingUserModifier

	if len(memories) < newUserMemoryCount {

		log.New("Using new user prompt modifier").AddUser(recipient).Log()
		promptModifier = NewUserModifier
	}

	if err != nil {
		log.New("Error remembering history %s", err.Error()).
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	pst, err := time.LoadLocation("America/Los_Angeles") // PST is often represented by the America/Los_Angeles timezone.

	if err != nil {
		log.New("Error loading PST location: %s. Exiting.", err.Error())

		return
	}

	// Convert date to PST.  In the future we will use the user's timezone
	pstDate := nowInUTC.In(pst)
	formattedDate := pstDate.Format("January 2, 2006 3:04pm")

	prompt := fmt.Sprintf(ConditioningPrompt, promptModifier, formattedDate, recipient.Firstname)

	log.New("Generated Prompt").Add("prompt", prompt).
		AddUser(recipient).AddMessage(&msg).
		Add("memory_count", strconv.Itoa(len(memories))).
		Log()

	// Send the prompt for completion
	completion, err := h.CompletionService.GetCompletion(msg.Body, prompt, &memories)

	// Strip some non GSM characters from the outbound message
	completion = h.CompletionService.CleanCompletionText(completion)

	if err != nil {
		log.New("Error getting completion").Add("prompt", prompt).
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).
			Add("memory_count", strconv.Itoa(len(memories))).Log()

		return
	}

	// Create a message entry in the db
	newMessage := NewMessage(&msg)
	newMessage.ConversationID = msg.ConversationID
	newMessage.FromUserID = models.GetSystemUser().ID
	newMessage.ToUserID = recipient.ID
	newMessage.SentAt = &nowInUTC
	newMessage.MessageStatus = models.NewMessageStatusSending()

	newMessage.Body = completion
	err = h.MessageService.CreateMessage(newMessage)

	if err != nil {
		log.New("Error saving new message").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	log.New("Sending SMS from %s to %s",
		newMessage.From.PhoneNumber,
		recipient.PhoneNumber).
		Add("message", newMessage.Body).
		Log()

	// Send the message back to the sender
	smsResponse, err := twilio.SendSMS(
		models.GetSystemUser().PhoneNumber, recipient.PhoneNumber, completion,
	)

	if err != nil {
		log.New("Error: Sending sms message").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	if smsResponse == nil {
		log.New("Error: SMS Response is empty").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	// Update the message with its SID from twilio
	newMessage.ReferenceID = smsResponse.Sid
	err = h.MessageService.UpdateMessage(newMessage)

	if err != nil {
		log.New("Error: Updating new message with reference ID").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	log.New("Successfully queued outbound message from %s to %s",
		msg.MessageType.Name, *smsResponse.To).
		AddSmsResponse(smsResponse).
		Log()

	defer func() {

		if r := recover(); r != nil {

			log.New("Panic while trying to process emotions: %v", r).Log()
		}

	}()
	h.ProcessEmotions(recipient, msg, event)

}

func (h *SendSMSLambdaHandler) ProcessEmotions(recipient *models.User, msg models.Message, event events.SQSMessage) {

	scores, err := h.NRCLexService.ProcessMessage(recipient, &msg)

	if err != nil {

		log.New("Error: Processing NRC lex message: %s", err.Error()).
			AddError(err).
			AddUser(recipient).
			AddSQSEvent(&event).
			AddError(err).
			AddMessage(&msg).
			Log()

		return
	}

	com := strconv.FormatFloat(scores.VaderEmotionScore.Compound, 'f', 4, 64)
	pos := strconv.FormatFloat(scores.VaderEmotionScore.Pos, 'f', 4, 64)
	neg := strconv.FormatFloat(scores.VaderEmotionScore.Neg, 'f', 4, 64)
	neu := strconv.FormatFloat(scores.VaderEmotionScore.Neu, 'f', 4, 64)

	emotionBlob, _ := json.Marshal(scores)

	log.New("Successfully processed NRC lex message for %s", recipient.PhoneNumber).
		Add("compound_sentiment", com).
		Add("neg_sentiment", neg).
		Add("neu_sentiment", neu).
		Add("positive_sentiment", pos).
		Add("raw_scores", string(emotionBlob)).
		AddUser(&msg.To).
		AddSQSEvent(&event).
		AddError(err).
		AddMessage(&msg).
		Log()
}

func (h *SendSMSLambdaHandler) GetMemories(recipient *models.User, event events.SQSMessage, msg models.Message) ([]models.Message, error) {

	lastFewMemories, err := h.MemoryService.GetLastNMessagePairs(recipient, h.MaxLastFewMemories)

	if err != nil {
		log.New("Error retrieving last few memories for user %s", recipient.PhoneNumber).
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return nil, err
	}

	aFewOlderMemories, err := h.MemoryService.GetRandomMessagePairs(recipient, h.MaxMemories)

	myOldMemories := utils.FilterSlice(*aFewOlderMemories, func(m models.Message) bool {
		return m.FromUserID != models.GetSystemUser().ID
	})

	if err != nil {
		log.New(
			"Error retrieving a few older memories for user %s",
			recipient.PhoneNumber).Log()

		return nil, err
	}

	memories := append(*lastFewMemories, myOldMemories...)
	slices.Reverse(memories)

	return memories, nil
}

func NewMessage(incomingMessage *models.Message) *models.Message {

	return &models.Message{
		FromUserID:      models.GetSystemUser().ID,
		ToUserID:        incomingMessage.FromUserID,
		ReferenceID:     incomingMessage.ReferenceID,
		MessageType:     models.NewMessageTypeSMS(),
		MessageStatusID: models.NewMessageStatusSent().ID,
	}
}

func ValidateEvent(sqsEvent *events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {

		return fmt.Errorf("no event records found")
	}

	if sqsEvent.Records[0].Body == "" {

		return fmt.Errorf("no event body found")
	}

	return nil
}

func main() {

	log.New("SMS Sender Lambda booting.....").Log()

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
		log.New("Error pinging Google. Possible bad internet connection.").AddError(err).Log()

		return
	}

	memoryService := message.NewMemoryService(
		message.NewMessageRepository(database),
	)

	nrclexRepo := &nrclex.Repository{DB: database}

	restClient := utils.NewRestClient()
	nrcClient := nrclex.NewNRCLexClient(restClient.GetClient())

	handler := &SendSMSLambdaHandler{

		MaxMemories:        maxMemories + maxLastFewMessages,
		MaxLastFewMemories: maxLastFewMessages,

		CompletionService: &ai.OpenAICompletionService{
			RemoveEmojis: false,
		},
		MemoryService: memoryService,
		NRCLexService: emotions.NewNRCLexService(nrcClient, nrclexRepo),
	}

	log.New("SMS Sender Lambda ready. Initializing.").Log()

	handler.Init(database)

	log.New("SMS Sender Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
