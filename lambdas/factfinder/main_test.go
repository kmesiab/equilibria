package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/facts"
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
	require.NoError(t, err, "Could not run tests, could not set up mock db")

	completionSvc := &ai.OpenAICompletionService{
		RemoveEmojis: false,
	}

	factRepo := facts.NewRepository(db)
	factSvc := facts.NewService(factRepo, completionSvc)

	handler := &FactFinderLambdaHandler{
		Service: factSvc,
	}

	handler.Init(db)
	err = handler.HandleRequest(event)
	assert.Error(t, err, "no event body found")
}

func TestHandleRequest(t *testing.T) {
	//test.SetEnvVars()
	event := events.SQSEvent{
		Records: []events.SQSMessage{{
			Body: `{
                "id": 1,
                "reference_id": "SAxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
                "conversation_id": 89,
                "from_user_id": 36,
                "to_user_id": 1,
                "body": "Today was a terrible day. My dog died and I'm so sad. Also I lowered my meds.'",
                "message_type_id": 2,
                "message_status_id": 3,
                "sent_at": null,
                "received_at": "2023-12-14T09:32:15Z",
                "created_at": "2023-12-14T09:32:15Z",
                "updated_at": "2023-12-14T09:32:15Z",
                "deleted_at": null,
                "conversation": {
                    "user": {
                        "status": {
                            "status_id": 2,
                            "status_name": "Active"
                        },
                        "id": 3,
                        "phone_number": "+12533243071",
                        "phone_verified": true,
                        "firstname": "New Name",
                        "lastname": "Mesiab",
                        "email": "kmesiab+equilibria_sms@gmail.com",
                        "account_status_id": 2
                    },
                    "id": 5,
                    "user_id": 3,
                    "start_time": "2023-12-14T09:32:15Z",
                    "end_time": null,
                    "created_at": "2023-12-14T09:32:15Z",
                    "updated_at": "0001-01-01T00:00:00Z"
                },
                "message_status": {
                    "status_id": 3,
                    "status_name": "Received"
                },
                "message_type": {
                    "id": 2,
                    "name": "SMS",
                    "bill_rate_in_credits": 1
                },
                "from": {
                    "status": {
                        "status_id": 2,
                        "status_name": "Active"
                    },
                    "id": 36,
                    "phone_number": "+12533243071",
                    "phone_verified": true,
                    "firstname": "New Name",
                    "lastname": "Mesiab",
                    "email": "kmesiab+equilibria_sms@gmail.com",
                    "account_status_id": 2
                },
                "to": {
                    "status": {
                        "status_id": 2,
                        "status_name": "Active"
                    },
                    "id": 1,
                    "phone_number": "+18333595081",
                    "phone_verified": true,
                    "firstname": "System",
                    "lastname": "User",
                    "email": "-@-",
                    "account_status_id": 2
                }
            }`,
		}},
	}

	var err error
	database := db.Get(config.Get())

	/*
		//db, mock, err := test.SetupMockDB()
		require.NoError(t, err, "Could not run tests, could not set up mock db")
		mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id`").
			WithArgs(3, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "nudge_enabled"}).
				AddRow(1, test.DefaultTestEmail, true))

	*/

	completionSvc := &ai.OpenAICompletionService{
		RemoveEmojis: false,
	}

	factRepo := facts.NewRepository(database)
	factSvc := facts.NewService(factRepo, completionSvc)

	handler := &FactFinderLambdaHandler{
		Service: factSvc,
	}

	handler.Init(database)
	err = handler.HandleRequest(event)
	assert.NoError(t, err)
}
