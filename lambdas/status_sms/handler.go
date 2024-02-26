package main

import (
	"net/http"
	"net/url"

	"github.com/Masterminds/formenc/encoding/form"
	"github.com/aws/aws-lambda-go/events"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// TwilioStatusEventHandler is a function that handles a specific status type.
type TwilioStatusEventHandler func(message *models.Message) error

// StatusSMSLambdaHandler handles Twilio status callbacks.
type StatusSMSLambdaHandler struct {
	lib.LambdaHandler
}

func (s *StatusSMSLambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !twilio.IsValidWebhookRequest(request, config.Get().TwilioAuthToken, false) {

		return log.New("Invalid webhook request signature. Rejecting webhook.").
			AddAPIProxyRequest(&request).
			Respond(http.StatusBadRequest)
	}

	var err error
	var vals url.Values

	log.New("Handling request").
		AddAPIProxyRequest(&request).Add("body", request.Body).
		Log()

	if vals, err = url.ParseQuery(request.Body); err != nil {

		return log.New("Error parsing form values").
			AddAPIProxyRequest(&request).AddError(err).Add("body", request.Body).
			Respond(http.StatusBadRequest)
	}

	var messageInfo = &models.TwilioMessageInfo{}

	if err = form.Unmarshal(vals, messageInfo); err != nil {

		return log.New("Error unmarshalling request").
			AddTwilioMessageInfo(messageInfo).
			AddAPIProxyRequest(&request).
			AddError(err).
			Add("body", request.Body).
			Respond(http.StatusBadRequest)
	}

	log.New("Status Update from %s to %s is currently %s",
		messageInfo.From, messageInfo.To, messageInfo.SMSStatus).
		Add("body", request.Body).
		AddAPIProxyRequest(&request).
		AddTwilioMessageInfo(messageInfo).Log()

	if err = s.ProcessMessage(messageInfo.GetTwilioMessageStatus(),
		messageInfo,
		s.DeductCredits,
		s.CloseConversation,
	); err != nil {

		return log.New("Error handling status for %s", messageInfo.MessageSid).
			AddAPIProxyRequest(&request).AddTwilioMessageInfo(messageInfo).AddError(err).
			Respond(http.StatusInternalServerError)
	}

	log.New("Successfully processed status for %s", messageInfo.MessageSid).
		AddAPIProxyRequest(&request).Add("body", request.Body).AddTwilioMessageInfo(messageInfo).
		Log()

	// Success
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "text/xml"},
		StatusCode: http.StatusOK,
	}, nil
}

func (s *StatusSMSLambdaHandler) ProcessMessage(
	status models.TwilioMessageStatus,
	messageInfo *models.TwilioMessageInfo,
	deductCredits TwilioStatusEventHandler,
	closeConversation TwilioStatusEventHandler,
) error {

	// Get the message
	msg, err := s.MessageService.FindByReferenceID(messageInfo.SmsSid)

	if err != nil {

		return err
	}

	// convert and update its status
	messageStatus := models.ConvertTwilioStatusToMessageStatus(status)
	msg.MessageStatusID = messageStatus.ID

	err = s.MessageService.UpdateStatus(msg)

	if err != nil {
		return err
	}

	log.New("Updated %s message status to %s in the database",
		*msg.ReferenceID, messageStatus.Name)

	// Handle only the case where we have a failure or success.
	switch status {

	// Success cases
	case models.TwilioMessageStatusDelivered:
		err = deductCredits(msg)

		if err != nil {

			return err
		}

		return closeConversation(msg)

		// Fail cases
	case models.TwilioMessageStatusFailed, models.TwilioMessageStatusUndelivered:

		return closeConversation(msg)
	}

	return nil
}

func (s *StatusSMSLambdaHandler) DeductCredits(msg *models.Message) error {
	log.New("Deducting credits from %s %s", msg.To.Firstname, msg.To.PhoneNumber).Log()
	return nil
}

func (s *StatusSMSLambdaHandler) CloseConversation(msg *models.Message) error {

	log.New("Closing the conversation %d...", msg.ConversationID).
		AddMessage(msg).Log()

	conversationID, err := s.MessageService.FindConversationIDByReferenceID(*msg.ReferenceID)

	if err != nil {
		log.New("Error getting conversation ID for %d", msg.ConversationID).
			AddMessage(msg).AddError(err).Log()

		return err
	}

	return s.ConversationService.EndConversation(conversationID)
}
