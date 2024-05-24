package report_generator_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/report_generator"
)

func TestHandleRequestMethodNotAllowed(t *testing.T) {
	// Setup mock database
	db, _, err := test.SetupMockDB()
	require.NoError(t, err, "Setting up mock db failed")

	memSvc := message.NewMemoryService(
		message.NewMessageRepository(db),
	)

	llmSvc := &ai.OpenAICompletionService{
		RemoveEmojis: false,
	}

	// Initialize the handler with mocked services as needed
	handler := report_generator.PatientReportGeneratorLambdaHandler{
		MemoryService:     memSvc,
		CompletionService: llmSvc,
	}

	handler.Init(db)

	// Create an API Gateway Proxy Request with an unsupported method
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
	}

	// Call the handler's HandleRequest method
	response, err := handler.HandleRequest(request)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode, "Expected a 405 Method Not Allowed response")
}
