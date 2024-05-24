package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/hasher"
	"github.com/kmesiab/equilibria/lambdas/lib/jwt"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type ManageUserLambdaHandler struct {
	lib.LambdaHandler
	KeyRotator   jwt.KeyRotatorInterface
	TokenService jwt.TokenServiceInterface
}

func (h *ManageUserLambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "POST":

		return h.Create(request)
	case "PUT":

		return h.Update(request)
	case "GET":

		return h.GetUser(request)

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

func (h *ManageUserLambdaHandler) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		newUser models.User
		err     error
	)

	err = json.Unmarshal([]byte(request.Body), &newUser)

	// Is the user valid?
	if err != nil || !newUser.IsValid() {

		return lib.RespondWithError("Missing required fields", err, http.StatusBadRequest)
	}

	// Is it a valid email address?
	if _, err = mail.ParseAddress(newUser.Email); err != nil {

		return lib.RespondWithError("Invalid email address", nil, http.StatusBadRequest)
	}

	// Is it a valid phone number?
	if newUser.PhoneNumber == "" || !twilio.IsValidPhoneNumber(newUser.PhoneNumber) {
		msg := fmt.Sprintf("Invalid phone number %s", newUser.PhoneNumber)

		return lib.RespondWithError(msg, nil, http.StatusBadRequest)
	}

	// Is the password secure?
	if *newUser.Password == "" || !hasher.IsSecureString(*newUser.Password) {

		return lib.RespondWithError(
			"Password must be at least 8 characters", nil, http.StatusBadRequest)
	}

	// Pending activation
	newUser.AccountStatusID = 1
	newUser.EnableNudges()

	hashedPassword, err := hasher.HashPassword(*newUser.Password)

	if err != nil {
		return lib.RespondWithError("Error hashing password", err, http.StatusInternalServerError)
	}

	newUser.Password = &hashedPassword

	// Create the user
	err = h.UserService.Create(&newUser)

	if err != nil {

		msg := "Error creating user"

		if db.IsDuplicateEntryError(err) {
			msg = "Phone number already registered"
		}

		return lib.RespondWithError(msg, nil, http.StatusBadRequest)
	}

	token, err := h.TokenService.Issue(&newUser, 10, h.KeyRotator)

	if err != nil {
		return lib.RespondWithError(
			"Error reissuing JWT", err, http.StatusInternalServerError)
	}

	userResponse := models.MakeUserResponseFromUser(&newUser)

	response := struct {
		User  *models.UserResponse `json:"user"`
		Token string               `json:"token"`
	}{
		User:  userResponse,
		Token: token,
	}

	responseBytes, err := json.Marshal(response)

	if err != nil {

		return lib.RespondWithError("Error marshaling newly created user", err, http.StatusBadRequest)
	}

	return events.APIGatewayProxyResponse{
		Headers:    config.DefaultHttpHeaders,
		StatusCode: http.StatusCreated,
		Body:       string(responseBytes)}, nil
}

func (h *ManageUserLambdaHandler) Update(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		inputUser   = &models.User{}
		updatedUser *models.User
		err         error
	)

	err = json.Unmarshal([]byte(request.Body), inputUser)

	if err != nil {

		return lib.RespondWithError("Invalid request body", err, http.StatusBadRequest)
	}

	if inputUser.ID == 0 {

		return lib.RespondWithError("Invalid user ID", nil, http.StatusBadRequest)
	}

	if inputUser.PhoneNumber != "" && !twilio.IsValidPhoneNumber(inputUser.PhoneNumber) {
		msg := fmt.Sprintf("Invalid phone number %s", inputUser.PhoneNumber)

		return lib.RespondWithError(msg, nil, http.StatusBadRequest)
	}

	if inputUser.Password != nil {
		err := h.updatePassword(inputUser)

		if err != nil {

			return lib.RespondWithError("Error updating password",
				err, http.StatusInternalServerError)
		}
	}

	log.New("Updating user %d", inputUser.ID).AddUser(inputUser).Log()

	err = h.UserService.Update(inputUser)

	if err != nil {

		msg := "Error updating user"

		if db.IsDuplicateEntryError(err) {
			msg = "Phone number already registered"
		}

		return lib.RespondWithError(msg, err, http.StatusBadRequest)
	}

	updatedUser, err = h.UserService.GetUserByID(inputUser.ID)

	if err != nil {

		return lib.RespondWithError("Couldn't get newly updated user",
			err, http.StatusInternalServerError)
	}

	token, err := h.TokenService.Issue(updatedUser, 10, h.KeyRotator)

	if err != nil {
		return lib.RespondWithError(
			"Error reissuing JWT", err, http.StatusInternalServerError)
	}

	msg := log.New("User has been updated and token refreshed").
		Add("token", token).
		Write()

	return events.APIGatewayProxyResponse{
		Headers:    config.DefaultHttpHeaders,
		StatusCode: http.StatusOK,
		Body:       msg,
	}, nil
}

func (h *ManageUserLambdaHandler) GetUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		intUserID int          // userID represented as an integer
		userID    string       // the path parameter
		newUser   *models.User // the user we fetched
		err       error        // any errors
		ok        bool         // success check
	)

	if userID, ok = request.PathParameters["userId"]; !ok {

		return lib.RespondWithError("User ID is required "+userID, nil, http.StatusBadRequest)
	}

	// Disambiguate the identifier type (phone number or numeric ID)
	if twilio.IsValidPhoneNumber(userID) {
		// If it's a phone number, parse and fetch it.
		newUser, err = h.UserService.GetUserByPhoneNumber(userID)
	} else {

		// Otherwise, get it by its numeric ID
		intUserID, err = strconv.Atoi(userID)

		if err != nil {

			return lib.RespondWithError("Invalid user ID", err, http.StatusBadRequest)
		}

		newUser, err = h.UserService.GetUserByID(int64(intUserID))
	}

	// If the user does not exist
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return log.New("User not found").
				Add("id", userID).
				Respond(http.StatusNotFound)
		}

		return log.New("Error retrieving user").
			AddError(err).
			Respond(http.StatusInternalServerError)

	}

	// If the user can't be fetched.  This will usually be caught above.
	if newUser == nil || newUser.ID == 0 {
		return log.New("User not found").
			Add("id", userID).
			Respond(http.StatusNotFound)
	}

	userResponse := models.MakeUserResponseFromUser(newUser)
	responseBytes, err := json.Marshal(userResponse)

	if err != nil {

		return log.New("Error marshaling user").
			AddError(err).
			Respond(http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Headers:    config.DefaultHttpHeaders,
		StatusCode: http.StatusOK,
		Body:       string(responseBytes),
	}, nil
}

func (h *ManageUserLambdaHandler) updatePassword(user *models.User) error {

	if *user.Password == "" || !hasher.IsSecureString(*user.Password) {

		return fmt.Errorf("password must be at least 8 characters")
	}

	hashedPassword, err := hasher.HashPassword(*user.Password)

	if err != nil {

		return fmt.Errorf("error hashing password\"")
	}

	user.Password = &hashedPassword

	return nil
}

func main() {
	log.New("Manage User Lambda booting...").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config").Log()
	}

	database := db.Get(cfg)
	handler := &ManageUserLambdaHandler{
		KeyRotator:   &jwt.KeyRotator{},
		TokenService: jwt.TokenService{},
	}
	handler.Init(database)

	log.New("Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
