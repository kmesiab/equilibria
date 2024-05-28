package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/facts"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/sqs"
	"github.com/kmesiab/equilibria/lambdas/lib/utils"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// FactFinderLambdaHandler handles the FactFinder Lambda logic
type FactFinderLambdaHandler struct {
	lib.LambdaHandler

	Service facts.ServiceInterface
}

func (h *FactFinderLambdaHandler) HandleRequest(sqsEvent events.SQSEvent) error {

	defer func() {
		if r := recover(); r != nil {
			log.New("Panic while processing event: %v\nStack trace:\n%s", r, debug.Stack()).Log()
		}
	}()

	if len(sqsEvent.Records) == 0 {
		log.New("No records found in the event.  Shutting down.").Log()

		return nil
	}

	for _, record := range sqsEvent.Records {
		err := h.processMessage(record)

		if err != nil {
			log.New("Error processing message: %s, error: %v", record.Body, err)
			return err
		}
	}

	return nil
}

func (h *FactFinderLambdaHandler) processMessage(record events.SQSMessage) error {
	var (
		err         error
		currentUser *models.User

		msg         models.Message
		eventRecord sqs.SQSEventRecord
	)

	if err = ValidateEvent(record); err != nil {
		log.New("Event record %s had no body", record.MessageId).
			AddMap(record.Attributes).
			AddError(err).Log()

		return err
	}

	// Unpack the SNS Event Record
	if err = json.Unmarshal([]byte(record.Body), &eventRecord); err != nil {
		log.New("Error unmarshalling event record").
			AddError(err).Log()

		return err
	}

	// Unpack the message from the event record
	if err = json.Unmarshal([]byte(eventRecord.Message), &msg); err != nil {
		log.New("Error unmarshalling message from event record").
			AddError(err).Log()

		return err
	}

	// Ignore system user
	if msg.FromUserID == models.GetSystemUser().ID {
		log.New("Ignoring message from system user").Log()

		return err
	}

	if currentUser, err = h.UserService.GetUserByID(msg.FromUserID); err != nil {
		log.New("Could not locate user %d.  Rejecting.", msg.FromUserID).
			AddError(err).AddMessage(&msg).Log()

		return err
	}

	log.New("Fact-finding request for: %s", currentUser.PhoneNumber).Log()

	identifiedFacts, err := h.Service.FindFacts(msg.Body)

	if err != nil {
		log.New("Error in FindFacts").AddError(err).Log()

		return err
	}

	// If we detected facts...
	if identifiedFacts != nil && len(*identifiedFacts) > 0 {
		for _, fact := range *identifiedFacts {

			// Map the fact
			f := &models.Fact{
				UserID:         currentUser.ID,
				ConversationID: msg.ConversationID,
				Body:           fact.Fact,
				Reasoning:      fact.Reasoning,
			}

			err := h.Service.CreateFact(f)

			if err != nil {
				log.New("Error saving fact: %s", fact.Fact).AddError(err).Log()

				return err
			}

		}
	} else {

		log.New("No facts identified for %s", currentUser.PhoneNumber).Log()
	}

	log.New("Fact-finding complete for: %s", currentUser.PhoneNumber).Log()
	return nil
}

func ValidateEvent(record events.SQSMessage) error {
	if record.Body == "" {
		return fmt.Errorf("no event body found")
	}

	return nil
}

func main() {
	log.New("FactFinder Lambda booting...").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config")
	}

	database := db.Get(cfg)

	if err := utils.PingDatabase(database); err != nil {
		log.New("Error pinging database").AddError(err).Log()

		return
	}

	completionSvc := &ai.OpenAICompletionService{
		RemoveEmojis: true,
	}

	factRepo := facts.NewRepository(database)
	factSvc := facts.NewService(factRepo, completionSvc)

	handler := &FactFinderLambdaHandler{
		Service: factSvc,
	}

	handler.Init(database)

	lambda.Start(handler.HandleRequest)
}
