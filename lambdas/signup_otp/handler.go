package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type SignupOTPLambdaHandler struct {
	lib.LambdaHandler
}

type OTPInputPayload struct {
	Code        string `form:"code" json:"code"`
	PhoneNumber string `form:"phone_number" json:"phone_number"`
}

func (s *SignupOTPLambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "POST":

		return s.Create(request)
	case "PUT":

		return s.Update(request)

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

func (s *SignupOTPLambdaHandler) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		user *models.User
		err  error
	)

	payload := &OTPInputPayload{}
	err = json.Unmarshal([]byte(request.Body), &payload)

	if err != nil {

		return lib.RespondWithError("Invalid json body", nil, http.StatusNotModified)
	}

	if payload.PhoneNumber == "" || !twilio.IsValidPhoneNumber(payload.PhoneNumber) {

		return log.New("Invalid phone number %s", payload.PhoneNumber).
			Respond(http.StatusBadRequest)

	}

	user, err = s.UserService.GetUserByPhoneNumber(payload.PhoneNumber)

	if err != nil {

		return lib.RespondWithError(
			"Couldn't find a user with this phone number", err, http.StatusBadRequest,
		)
	}

	if user.PhoneVerified {

		log.New("User %s already verified", user.PhoneNumber).
			Log()
	}

	log.New("Sending OTP Request to Twilio").
		Add("phone_number", user.PhoneNumber).Log()

	signupOtpResponse, err := twilio.SendOTP(user.PhoneNumber)

	if err != nil {

		return lib.RespondWithError("Couldn't send OTP", err, http.StatusInternalServerError)
	}

	// Success
	return log.New("Sent OTP to %s", user.PhoneNumber).
		Add("status", *signupOtpResponse.Status).Respond(http.StatusOK)

}

func (s *SignupOTPLambdaHandler) Update(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.New("Body : %s", request.Body).Log()

	payload := &OTPInputPayload{}
	err := json.Unmarshal([]byte(request.Body), &payload)

	if err != nil {

		return lib.RespondWithError("Invalid json body", nil, http.StatusNotModified)
	}

	if payload.PhoneNumber == "" || !twilio.IsValidPhoneNumber(payload.PhoneNumber) {

		return log.New("Invalid phone number %s", payload.PhoneNumber).
			Respond(http.StatusBadRequest)

	}

	if payload.Code == "" {

		return lib.RespondWithError("A valid code is required", nil, http.StatusBadRequest)
	}

	user, err := s.UserService.GetUserByPhoneNumber(payload.PhoneNumber)

	if err != nil {

		return lib.RespondWithError(
			"Couldn't find an active user with this phone number", err, http.StatusBadRequest,
		)
	}

	if user.PhoneVerified {
		log.New("User %s already verified", user.PhoneNumber).
			Log()
	}

	log.New("Verifying OTP").Add("phone_number", user.PhoneNumber).Log()

	signupOtpResponse, err := twilio.VerifyOTP(payload.PhoneNumber, payload.Code)

	if err != nil {

		log.New("Error verifying OTP").AddError(err).
			Add("phone_number", payload.PhoneNumber).Log()

		return lib.RespondWithError("", nil, http.StatusNotModified)
	}

	if *signupOtpResponse.Status != "approved" {
		return log.New("Incorrect code for %s", payload.PhoneNumber).
			Add("phone_number", *signupOtpResponse.To).
			Add("status", *signupOtpResponse.Status).
			Add("channel", *signupOtpResponse.Channel).
			Respond(http.StatusBadRequest)
	}

	user.PhoneVerified = true
	user.AccountStatusID = models.AccountStatusActive

	if err := s.UserService.Update(user); err != nil {

		return log.New("Couldn't update user").AddError(err).
			Respond(http.StatusInternalServerError)
	}

	// Success
	return log.New("OTP Approved for %s", payload.PhoneNumber).
		Add("status", *signupOtpResponse.Status).Respond(http.StatusOK)
}
