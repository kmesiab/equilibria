package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type MockSQSSender struct{}

func (m MockSQSSender) Send(_ string, _ *models.Message) error {
	return nil
}

func TestReceiveSMSLambdaHandler_Receive(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could not set up mock db")

	handler := ReceiveSMSLambdaHandler{
		SQSSender: MockSQSSender{},
	}
	handler.Init(db)

	// First we look up the user who sent the message
	test.ExpectMockSelectUser(&mock, "+12533243071")

	// We use the user and the data to create a new conversation
	test.ExpectMockInsertConversation(&mock)

	// Then we add the message and attach it to the conversation
	test.ExpectMockInsertMessage(&mock)

	// Then we look up the message
	mock.ExpectQuery("SELECT \\* FROM `messages` WHERE id").
		WithArgs(1).
		WillReturnRows(test.GenerateMockMessageRepositoryMessages())

	// And the conversation
	mock.ExpectQuery("SELECT \\* FROM `conversations`").WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))
	// And the user associated with each
	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectUser(&mock, 1)

	mock.ExpectQuery("SELECT \\* FROM `message_statuses`").
		WithArgs(models.NewMessageStatusPending().ID).WillReturnRows(test.GenerateMockMessageStatus())

	mock.ExpectQuery("SELECT \\* FROM `message_types`").
		WithArgs(models.NewMessageTypeSMS().ID).WillReturnRows(test.GenerateMockMessageType())

	test.ExpectMockSelectUser(&mock, 1)

	response, err := handler.Receive(events.APIGatewayProxyRequest{
		Body: "ToCountry=US&ToState=&SmsMessageSid=SM62876cd3611d64defdece80d9aa1f703&NumMedia=0&ToCity=&FromZip=98106&SmsSid=SM62876cd3611d64defdece80d9aa1f703&FromState=WA&SmsStatus=received&FromCity=SEATTLE&Body=Why+do+you+say+that%3F+&FromCountry=US&To=%2B18333595081&MessagingServiceSid=MGa3799c565299f143097ff388571be2b2&ToZip=&NumSegments=1&MessageSid=SM62876cd3611d64defdece80d9aa1f703&AccountSid=AC0e8e16c274b3ae7740b1a854b8c9846a&From=%2B12533243071&ApiVersion=2010-04-01",
	})

	require.NoError(t, err,
		"error should be nil when handling sms receive request")

	assert.Equal(t, 201, response.StatusCode, "Receive SMS should return 200")

}

func TestHandleRequest_ValidTwilioSignature(t *testing.T) {

	test.SetEnvVars()

	db, mock, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	test.ExpectMockSelectUser(&mock, "+12533243071")
	test.ExpectMockInsertConversation(&mock)

	// Then we add the message and attach it to the conversation
	test.ExpectMockInsertMessage(&mock)

	mock.ExpectQuery("SELECT \\* FROM `messages` WHERE id").
		WithArgs(1).
		WillReturnRows(test.GenerateMockMessageRepositoryMessages())

	// And the conversation
	mock.ExpectQuery("SELECT \\* FROM `conversations`").WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	// And the user associated with each
	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectUser(&mock, 1)

	// Then fetch the status and type of message
	mock.ExpectQuery("SELECT \\* FROM `message_statuses`").
		WithArgs(models.NewMessageStatusPending().ID).
		WillReturnRows(test.GenerateMockMessageStatus())

	mock.ExpectQuery("SELECT \\* FROM `message_types`").
		WithArgs(models.NewMessageTypeSMS().ID).
		WillReturnRows(test.GenerateMockMessageType())

	test.ExpectMockSelectUser(&mock, 1)

	mock.ExpectQuery("SELECT \\* FROM `messages` WHERE id").
		WithArgs(1).WillReturnRows(
		test.GenerateMockMessageRepositoryMessages(),
	)

	handler := ReceiveSMSLambdaHandler{
		SQSSender: MockSQSSender{},
	}
	handler.Init(db)

	req := *test.GenerateMockTwilioFormPostRequest()
	req.HTTPMethod = "INVALID HTTP METHOD"

	response, err := handler.HandleRequest(req)

	require.NoError(t, err,
		"error should be nil when handling sms receive request")

	assert.Equal(t, 405, response.StatusCode,
		"Receive SMS should return Invalid method")

}

func TestHandleRequest_InvalidTwilioSignature(t *testing.T) {

	test.SetEnvVars()

	db, _, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	handler := ReceiveSMSLambdaHandler{}
	handler.Init(db)

	response, err := handler.HandleRequest(events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"x-twilio-signature": "invalid-signature",
			"Host":               "localhost:3000",
			"Path":               "/sms/receive",
			"Content-Type":       twilio.WebhookContentType,
		},
		Body: "Body=sample%20message%3F&From=%2B12533243071&To=18333595081&MessageSid=SAxx&AccountSid=ACxx",
	})

	require.NoError(t, err,
		"error should be nil when handling sms receive request")

	assert.Equal(t, 400, response.StatusCode,
		"An invalid Twilio signature should return 400")

}
