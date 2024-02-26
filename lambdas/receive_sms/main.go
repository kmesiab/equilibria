// https://www.twilio.com/docs/messaging/guides/webhook-request#status-callback-parameters
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/form_unsmarshaler"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/sqs"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type ReceiveSMSLambdaHandler struct {
	lib.LambdaHandler
	SQSSender sqs.SenderInterface
}

func (h *ReceiveSMSLambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !twilio.IsValidWebhookRequest(request, config.Get().TwilioAuthToken, false) {

		return log.New("Invalid webhook request signature. Rejecting webhook.").
			AddAPIProxyRequest(&request).
			Respond(http.StatusBadRequest)
	}

	switch request.HTTPMethod {
	case "POST":

		return h.Receive(request)
	default:

		return log.New("Method %s not allowed", request.HTTPMethod).
			Respond(http.StatusMethodNotAllowed)
	}
}

func (h *ReceiveSMSLambdaHandler) Receive(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var (
		err     error
		message *models.Message // The newly created message object to be saved in the database
	)

	var sms = &models.TwilioMessageInfo{}

	err = form_unsmarshaler.UnMarshalBody(request, sms)

	if err != nil {

		return log.New("Error unmarshalling request body").
			AddError(err).Respond(http.StatusBadRequest)
	}

	if !twilio.IsValidPhoneNumber(sms.From) {

		return log.New("Invalid from phone number %s", sms.From).
			AddTwilioMessageInfo(sms).Respond(http.StatusBadRequest)
	}

	// Validate the sms
	if sms.From == "" || sms.To == "" || sms.Body == "" {

		return log.New("SMS missing required fields").
			AddTwilioMessageInfo(sms).Respond(http.StatusBadRequest)
	}

	// Log to cloudwatch
	log.New("SMS received. Starting conversation").
		AddTwilioMessageInfo(sms).Log()

	// Add this message to a new conversation
	message, err = h.StartConversation(sms)

	log.New("Conversation started, sending message to queue").Log()

	if err != nil {

		return log.New("Error starting a conversation for %s to %s", sms.From, sms.To).
			AddError(err).Respond(http.StatusInternalServerError)
	}

	// Send the message to the queue
	if err = h.SQSSender.Send(config.Get().SMSQueueURL, message); err != nil {

		fErr := h.Fail(message)
		return log.New("Error queueing message ID: %d\n", message.ID).
			AddError(err).AddError(fErr).Respond(http.StatusOK)
	}

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "text/xml"},
		StatusCode: http.StatusCreated,
	}, nil

}

func (h *ReceiveSMSLambdaHandler) Fail(message *models.Message) error {

	now := time.Now()
	message.MessageStatusID = models.NewMessageStatusFailed().ID
	message.Conversation.EndTime = &now

	if err := h.ConversationService.UpdateConversation(&message.Conversation); err != nil {

		return err
	}

	return nil
}

func (h *ReceiveSMSLambdaHandler) StartConversation(sms *models.TwilioMessageInfo) (*models.Message, error) {

	var (
		fromUser     *models.User         // The user identified by phone number
		toUser       *models.User         // The system user
		msg          *models.Message      // Create a new message from this sms
		conversation *models.Conversation // Create a new conversation for the message
		err          error
	)

	toUser = models.GetSystemUser()

	if fromUser, err = h.UserService.GetUserByPhoneNumber(sms.From); err != nil {

		return nil, fmt.Errorf("error getting user %s: %s", sms.From, err)
	}

	// Package the sms into a message struct
	msg = h.NewMessage(sms, fromUser, toUser)

	// Every inbound message gets a new conversation
	if conversation, err = h.CreateConversation(msg); err != nil {

		return nil, fmt.Errorf("error creating conversation: %s", err)
	}

	// Link the conversation and message
	msg.ConversationID = conversation.ID

	// Create the message in the database
	if err = h.MessageService.CreateMessage(msg); err != nil {

		return nil, fmt.Errorf("error creating message: %s", err)
	}

	// Fetch the newly created message
	if msg, err = h.MessageService.FindByID(msg.ID); err != nil {

		return nil, fmt.Errorf("error fetching newly created message: %s", err)
	}

	return msg, nil

}

func (h *ReceiveSMSLambdaHandler) CreateConversation(msg *models.Message) (*models.Conversation, error) {

	now := time.Now()
	conversation := &models.Conversation{
		UserID:    msg.FromUserID,
		StartTime: &now,
	}

	err := h.ConversationService.CreateConversation(conversation)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (h *ReceiveSMSLambdaHandler) NewMessage(sms *models.TwilioMessageInfo, fromUser, toUser *models.User) *models.Message {
	now := time.Now()
	msg := &models.Message{
		ReferenceID:     &sms.SmsSid,
		FromUserID:      fromUser.ID,
		ToUserID:        toUser.ID,
		ReceivedAt:      &now,
		MessageTypeID:   models.NewMessageTypeSMS().ID,
		MessageStatusID: models.NewMessageStatusReceived().ID,
	}
	msg.Body = sms.Body

	return msg
}

func main() {

	log.New("Receive Lambda booting...").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config")
	}

	database := db.Get(cfg)

	handler := ReceiveSMSLambdaHandler{
		SQSSender: sqs.NewSQSSender(),
	}
	handler.Init(database)

	log.New("Receive Lambda invoking...").Log()

	lambda.Start(handler.HandleRequest)
}
