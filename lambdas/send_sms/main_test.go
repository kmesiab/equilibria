package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

func TestHandleRequest_NoBody(t *testing.T) {
	test.SetEnvVars()
	event := events.SQSEvent{
		Records: []events.SQSMessage{{
			Body: "",
		}},
	}

	db, _, err := test.SetupMockDB()

	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	handler := &SendSMSLambdaHandler{
		CompletionService: &ai.OpenAICompletionService{
			RemoveEmojis: false,
		},
		MemoryService: &message.MemoryService{
			MessageService: message.NewMessageService(message.NewMessageRepository(db)),
		},
	}

	handler.Init(db)
	handler.HandleRequest(event)
}

func TestHandleRequest(t *testing.T) {
	test.SetEnvVars()
	event := events.SQSEvent{
		Records: []events.SQSMessage{{
			Body: "{\r\n\t\"id\": 1,\r\n\t\"reference_id\": \"SAxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\r\n\t\"conversation_id\": 5,\r\n\t\"from_user_id\": 3,\r\n\t\"to_user_id\": 1,\r\n\t\"body\": \"Hey, I was wondering if you thought maybe i'd be getting that promotion?\",\r\n\t\"message_type_id\": 2,\r\n\t\"message_status_id\": 3,\r\n\t\"sent_at\": null,\r\n\t\"received_at\": \"2023-12-14T09:32:15Z\",\r\n\t\"created_at\": \"2023-12-14T09:32:15Z\",\r\n\t\"updated_at\": \"2023-12-14T09:32:15Z\",\r\n\t\"deleted_at\": null,\r\n\t\"conversation\": {\r\n\t\t\"user\": {\r\n\t\t\t\"status\": {\r\n\t\t\t\t\"status_id\": 2,\r\n\t\t\t\t\"status_name\": \"Active\"\r\n\t\t\t},\r\n\t\t\t\"id\": 3,\r\n\t\t\t\"phone_number\": \"+12533243071\",\r\n\t\t\t\"phone_verified\": true,\r\n\t\t\t\"firstname\": \"New Name\",\r\n\t\t\t\"lastname\": \"Mesiab\",\r\n\t\t\t\"email\": \"kmesiab+equilibria_sms@gmail.com\",\r\n\t\t\t\"account_status_id\": 2\r\n\t\t},\r\n\t\t\"id\": 5,\r\n\t\t\"user_id\": 3,\r\n\t\t\"start_time\": \"2023-12-14T09:32:15Z\",\r\n\t\t\"end_time\": null,\r\n\t\t\"created_at\": \"2023-12-14T09:32:15Z\",\r\n\t\t\"updated_at\": \"0001-01-01T00:00:00Z\"\r\n\t},\r\n\t\"message_status\": {\r\n\t\t\"status_id\": 3,\r\n\t\t\"status_name\": \"Received\"\r\n\t},\r\n\t\"message_type\": {\r\n\t\t\"id\": 2,\r\n\t\t\"name\": \"SMS\",\r\n\t\t\"bill_rate_in_credits\": 1\r\n\t},\r\n\t\"from\": {\r\n\t\t\"status\": {\r\n\t\t\t\"status_id\": 2,\r\n\t\t\t\"status_name\": \"Active\"\r\n\t\t},\r\n\t\t\"id\": 5,\r\n\t\t\"phone_number\": \"+12533243071\",\r\n\t\t\"phone_verified\": true,\r\n\t\t\"firstname\": \"New Name\",\r\n\t\t\"lastname\": \"Mesiab\",\r\n\t\t\"email\": \"kmesiab+equilibria_sms@gmail.com\",\r\n\t\t\"account_status_id\": 2\r\n\t},\r\n\t\"to\": {\r\n\t\t\"status\": {\r\n\t\t\t\"status_id\": 2,\r\n\t\t\t\"status_name\": \"Active\"\r\n\t\t},\r\n\t\t\"id\": 1,\r\n\t\t\"phone_number\": \"+18333595081\",\r\n\t\t\"phone_verified\": true,\r\n\t\t\"firstname\": \"System\",\r\n\t\t\"lastname\": \"User\",\r\n\t\t\"email\": \"-@-\",\r\n\t\t\"account_status_id\": 2\r\n\t}\r\n}",
		}},
	}

	db, mock, err := test.SetupMockDB()

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id`").
		WithArgs(3, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "nudge_enabled"}).
			AddRow(1, test.DefaultTestEmail, true))

	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	handler := &SendSMSLambdaHandler{
		CompletionService: &ai.OpenAICompletionService{
			RemoveEmojis: false,
		},
		MemoryService: &message.MemoryService{
			MessageService: message.NewMessageService(message.NewMessageRepository(db)),
		},
	}

	handler.Init(db)
	handler.HandleRequest(event)
}
