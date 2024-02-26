package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/hasher"
	"github.com/kmesiab/equilibria/lambdas/lib/jwt"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type LoginLambda struct {
	lib.LambdaHandler
	KeyRotator   jwt.KeyRotatorInterface
	TokenService jwt.TokenServiceInterface
}

type LoginPayload struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (l *LoginLambda) HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "POST":

		return l.Login(request)

		// Enable cors Preflight
	case "OPTIONS":
		return events.APIGatewayProxyResponse{
			Headers:    config.DefaultHttpHeaders,
			StatusCode: http.StatusOK,
		}, nil

	default:

		return log.New("Method %s not allowed", request.HTTPMethod).
			Respond(http.StatusMethodNotAllowed)
	}
}

func (l *LoginLambda) Login(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		err          error
		user         *models.User
		loginPayload = &LoginPayload{}
	)

	// Get the body
	if err = json.Unmarshal([]byte(request.Body), loginPayload); err != nil {

		return log.New("Phone number and password required").
			AddError(err).Respond(http.StatusBadRequest)
	}

	log.New("Getting user with phone number %s", loginPayload.PhoneNumber).Log()

	// Get the user by phone number
	if user, err = l.UserService.GetUserByPhoneNumber(loginPayload.PhoneNumber); err != nil {

		return log.New("Invalid phone number or password").
			AddError(err).Respond(http.StatusBadRequest)
	}

	if !hasher.CheckPassword(loginPayload.Password, *user.Password) {

		return log.New("Password is incorrect").Respond(http.StatusBadRequest)
	}

	// Check the account status
	if user.AccountStatusID != models.AccountStatusActive {

		return log.New("User account is not active.").Respond(http.StatusBadRequest)
	}

	// Check the phone verification
	if !user.PhoneVerified {

		return log.New("User phone number is not verified.").Respond(http.StatusBadRequest)
	}

	token, err := l.TokenService.Issue(user, 10, l.KeyRotator)

	if err != nil {
		return lib.RespondWithError(
			"Error issuing JWT", err, http.StatusInternalServerError)
	}
	responseUser := models.MakeUserResponseFromUser(user)

	successResponse := struct {
		User *models.UserResponse `json:"user"`
		JWT  string               `json:"token"`
	}{
		User: responseUser,
		JWT:  token,
	}

	responseBytes, err := json.Marshal(successResponse)

	if err != nil {

		return log.New("Couldn't marshal user after verifying password.").
			AddError(err).Respond(http.StatusInternalServerError)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           config.DefaultHttpHeaders,
		MultiValueHeaders: nil,
		Body:              string(responseBytes),
	}

	response.Headers["Authorization"] = "Bearer: " + token

	return response, nil
}

func main() {

	log.New("Login Lambda booting...").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config")
	}

	handler := &LoginLambda{
		KeyRotator:   &jwt.KeyRotator{},
		TokenService: &jwt.TokenService{},
	}
	handler.Init(db.Get(cfg))

	log.New("Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
