package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/hasher"
	"github.com/kmesiab/equilibria/lambdas/lib/jwt"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

func TestLoginLambda_HandleRequest(t *testing.T) {

	test.SetEnvVars()

	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	loginLambda := &LoginLambda{
		KeyRotator:   jwt.NewMockKeyRotator(),
		TokenService: &jwt.TokenService{},
	}
	loginLambda.Init(db)

	// Mock user data
	phoneNumber := "+12533243071"
	password := "testPassword"

	hashedPassword, err := hasher.HashPassword(password)
	require.NoError(t, err, "Failed to hash password")

	var columnHeaders = []string{"phone_number", "password", "account_status_id", "phone_verified"}
	// Mock the database query for user retrieval
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE phone_number =").
		WithArgs(phoneNumber).
		WillReturnRows(sqlmock.NewRows(columnHeaders).
			AddRow(phoneNumber, hashedPassword, 2, true))

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs(2).WillReturnRows(test.GenerateMockAccountStatusActive())

	// Create request
	loginPayload := LoginPayload{PhoneNumber: phoneNumber, Password: password}
	requestBody, _ := json.Marshal(loginPayload)

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       string(requestBody),
	}

	// Test the HandleRequest function
	response, err := loginLambda.Login(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
