package main_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/jwt"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/manage_user"

	"github.com/kmesiab/equilibria/lambdas/models"
)

func TestManageUser_Post(t *testing.T) {

	test.SetEnvVars()

	randomPhoneNumber := fmt.Sprintf("+1253%d%d", rand.Intn(500), rand.Intn(1000))

	pwd := test.DefaultTestPassword

	user := models.User{
		Password:    &pwd,
		Firstname:   test.DefaultTestUserFirstname,
		Lastname:    test.DefaultTestUserLastname,
		PhoneNumber: randomPhoneNumber,
		UserTypeID:  1,
		Email:       test.DefaultTestEmail,
	}

	user.EnableNudges()

	userBytes, _ := json.Marshal(user)
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       string(userBytes),
	}

	db, mock, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WithArgs(
		sqlmock.AnyArg(),
		user.PhoneNumber,
		false,
		user.Firstname,
		user.Lastname,
		user.Email,
		1,
		1,
		user.NudgesEnabled(),
		user.ProviderCode,
	).
		WillReturnResult(
			test.GenerateMockLastAffectedRow(),
		)

	mock.ExpectCommit()

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs("Pending Activation", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Pending Activation"))

	handler := main.ManageUserLambdaHandler{
		TokenService: &jwt.TokenService{},
		KeyRotator:   jwt.NewMockKeyRotator(),
	}

	handler.Init(db)
	response, err := handler.Create(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	createdUserResponse := &struct {
		User  models.UserResponse `json:"user"`
		Token string              `json:"token"`
	}{}

	err = json.Unmarshal([]byte(response.Body), &createdUserResponse)

	assert.NoError(t, err)
	assert.Equal(t, user.Firstname, createdUserResponse.User.Firstname, "User should be set")
	assert.NotEmpty(t, createdUserResponse.Token, "Token should be set")
}

func TestManageUser_Update(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	user := models.User{
		ID:           3,
		Firstname:    "New Name",
		ProviderCode: "ABC123",
	}

	user.EnableNudges()

	userBytes, _ := json.Marshal(user)
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "PUT",
		Body:       string(userBytes),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `id`=\\?,`firstname`=\\?,`nudge_enabled`=\\?,`provider_code`=\\? WHERE id = \\?").WithArgs(
		user.ID, user.Firstname, user.NudgesEnabled(), user.ProviderCode, user.ID,
	).
		WillReturnResult(
			test.GenerateMockLastAffectedRow(),
		)

	mock.ExpectCommit()

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id`").
		WithArgs(user.ID, sqlmock.AnyArg()).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Pending Activation"))

	handler := main.ManageUserLambdaHandler{
		TokenService: &jwt.TokenService{},
		KeyRotator:   jwt.NewMockKeyRotator(),
	}
	handler.Init(db)
	response, err := handler.Update(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestManageUser_Get(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	// Mock the database queries
	userID := "1"

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(test.GenerateMockUserRepositoryUser()).
		WithArgs(int64(1), sqlmock.AnyArg())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	// Test the getUser function
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"userId": userID,
		},
	}

	handler := main.ManageUserLambdaHandler{
		KeyRotator: jwt.NewMockKeyRotator(),
	}
	handler.Init(db)
	response, err := handler.GetUser(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var returnedUser models.UserResponse
	err = json.Unmarshal([]byte(response.Body), &returnedUser)
	assert.NoError(t, err)
	assert.Equal(t, "jane", returnedUser.Firstname)
	assert.Equal(t, "doe", returnedUser.Lastname)
	assert.Equal(t, "2533243071", returnedUser.PhoneNumber)
	assert.Equal(t, "janedoe@email.com", returnedUser.Email)
}

func TestManageUser_Get404(t *testing.T) {
	test.SetEnvVars()

	db, mock, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	userID := "1"
	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(sqlmock.NewRows([]string{})).
		WithArgs(int64(1), sqlmock.AnyArg())

	// Test the getUser function
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"userId": userID,
		},
	}

	handler := main.ManageUserLambdaHandler{
		KeyRotator: jwt.NewMockKeyRotator(),
	}
	handler.Init(db)
	response, err := handler.GetUser(request)

	require.Nil(t, err, "Error should be nil when getting user by ID")
	assert.Equal(t, http.StatusNotFound, response.StatusCode,
		"Status code should be 404 when a user is not found")
}

func TestManageUser_InvalidCreateJson(t *testing.T) {
	db, _, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	// Test the getUser function
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Body:       "invalid json",
	}

	handler := main.ManageUserLambdaHandler{
		KeyRotator: jwt.NewMockKeyRotator(),
	}
	handler.Init(db)
	response, err := handler.Create(request)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	var responseErr = &test.JsonError{}
	err = json.Unmarshal([]byte(response.Body), responseErr)

	assert.NoError(t, err)
	assert.Equal(t, "Missing required fields", responseErr.Message)
}

func TestManageUser_InvalidUpdateJson(t *testing.T) {

	// Test the getUser function
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Body:       "invalid json",
	}

	db, _, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	handler := main.ManageUserLambdaHandler{
		KeyRotator: jwt.NewMockKeyRotator(),
	}
	handler.Init(db)
	response, err := handler.Update(request)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	var responseErr = &test.JsonError{}
	err = json.Unmarshal([]byte(response.Body), responseErr)

	assert.NoError(t, err)
	assert.Equal(t, "Invalid request body", responseErr.Message)
}

func TestManageUser_UpdateWithNoUserID(t *testing.T) {

	u := models.User{
		ID: 0,
	}

	bodyBytes, err := json.Marshal(u)
	require.NoError(t, err)

	// Test the getUser function
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Body:       string(bodyBytes),
	}

	db, _, err := test.SetupMockDB()
	require.NoError(t, err, "Could not run tests, could nto set up mock db")

	handler := main.ManageUserLambdaHandler{
		KeyRotator: jwt.NewMockKeyRotator(),
	}
	handler.Init(db)
	response, err := handler.Update(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	var responseErr = &test.JsonError{}
	err = json.Unmarshal([]byte(response.Body), responseErr)

	assert.NoError(t, err)
	assert.Equal(t, "Invalid user ID", responseErr.Message)
}
