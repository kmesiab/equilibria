package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

func TestMain_HandleRequest(t *testing.T) {

	test.SetEnvVars()
	request := test.GenerateMockTwilioFormPostRequest()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	mock.ExpectQuery("SELECT \\* FROM `messages` WHERE reference_id = ?").
		WithArgs("SMa74e33ba8361485b4bfbb6ec285ceac5").
		WillReturnRows(test.GenerateMockMessageRepositoryMessages())

	mock.ExpectQuery("SELECT \\* FROM `conversations`").WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectMessageStatusAndTypes(&mock)
	test.ExpectMockSelectUser(&mock, 1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `messages` SET `message_status_id`=").
		WithArgs(int64(3), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	test.ExpectMockSelectUser(&mock, 1)

	handler := &StatusSMSLambdaHandler{}
	handler.Init(db)

	response, err := handler.HandleRequest(*request)
	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Equal(t, "", response.Body)
}
