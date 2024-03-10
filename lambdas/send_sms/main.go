package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/lib/utils"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// How many past memories to include in the prompt
const maxMemories = 140

// How many immediately previous messages to include in the prompt
const maxLastFewMessages = 6

// How many memories you have to have before we consider you an 'existing'
// user, so the model treats you like it knows you well.
const newUserMemoryCount = 30

type SendSMSLambdaHandler struct {
	lib.LambdaHandler
	MemoryService     *message.MemoryService
	CompletionService ai.CompletionServiceInterface
}

func (h *SendSMSLambdaHandler) HandleRequest(sqsEvent events.SQSEvent) {

	var (
		err error
		now = time.Now()
		msg models.Message

		formattedDate = now.Format("January 2, 2006 3:04pm")
	)

	if err = ValidateEvent(&sqsEvent); err != nil {
		log.New("Error validating the event").AddError(err).Log()

		return
	}

	var event = sqsEvent.Records[0]
	var body = event.Body

	if err = json.Unmarshal([]byte(body), &msg); err != nil {
		log.New("Error unmarshalling message").AddError(err).Log()

		return
	}

	// Get the sender's user account
	recipient, err := h.UserService.GetUserByID(msg.FromUserID)

	if err != nil {
		log.New("Could not locate user %d.  Rejecting.", msg.FromUserID).
			AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	// Make sure the sender is not the system user
	if recipient.ID == models.GetSystemUser().ID {
		log.New("Message is from system user. Aborting.").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	log.New("Preparing a response for %s", recipient.PhoneNumber).
		AddUser(recipient).AddSQSEvent(&event).AddMessage(&msg).Log()

	// Get the memories for the user
	memories, err := h.GetMemories(recipient, event, msg)

	var promptModifier = ExistingUserModifier

	if len(memories) < newUserMemoryCount {
		promptModifier = NewUserModifier
	}

	if err != nil {
		log.New("Error remembering history %s", err.Error()).
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	prompt := fmt.Sprintf(ConditioningPrompt, promptModifier, formattedDate, recipient.Firstname)

	log.New("Generated Prompt").Add("prompt", prompt).
		AddUser(recipient).AddSQSEvent(&event).AddMessage(&msg).
		Add("memory_count", strconv.Itoa(len(memories))).Log()

	// Send the prompt for completion
	completion, err := h.CompletionService.GetCompletion(msg.Body, prompt, &memories)

	if err != nil {
		log.New("Error getting completion").Add("prompt", prompt).
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).
			Add("memory_count", strconv.Itoa(len(memories))).Log()

		return
	}

	// Create a message entry in the db
	newMessage := NewMessage(&msg)
	newMessage.ConversationID = msg.ConversationID
	newMessage.FromUserID = models.GetSystemUser().ID
	newMessage.ToUserID = recipient.ID
	newMessage.SentAt = &now
	newMessage.MessageStatus = models.NewMessageStatusSending()

	newMessage.Body = completion
	err = h.MessageService.CreateMessage(newMessage)

	if err != nil {
		log.New("Error saving new message").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	log.New("Sending SMS from %s to %s",
		newMessage.From.PhoneNumber,
		recipient.PhoneNumber).
		Add("message", newMessage.Body).
		Log()

	// Send the message back to the sender
	smsResponse, err := twilio.SendSMS(
		models.GetSystemUser().PhoneNumber, recipient.PhoneNumber, completion,
	)

	if err != nil {
		log.New("Error: Sending sms message").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	if smsResponse == nil {
		log.New("Error: SMS Response is empty").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	// Update the message with its SID from twilio
	newMessage.ReferenceID = smsResponse.Sid
	err = h.MessageService.UpdateMessage(newMessage)

	if err != nil {
		log.New("Error: Updating new message with reference ID").
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return
	}

	log.New("Successfully queued outbound message from %s to %s",
		msg.MessageType.Name, *smsResponse.To).
		AddSmsResponse(smsResponse).
		Log()
}

func (h *SendSMSLambdaHandler) GetMemories(recipient *models.User, event events.SQSMessage, msg models.Message) ([]models.Message, error) {

	lastFewMemories, err := h.MemoryService.GetLastNMessagePairs(recipient, maxLastFewMessages)

	if err != nil {
		log.New("Error retrieving last few memories for user %s", recipient.PhoneNumber).
			AddUser(recipient).AddSQSEvent(&event).AddError(err).AddMessage(&msg).Log()

		return nil, err
	}

	aFewOlderMemories, err := h.MemoryService.GetRandomMessagePairs(recipient, maxMemories)

	if err != nil {
		log.New(
			"Error retrieving a few older memories for user %s",
			recipient.PhoneNumber).Log()

		return nil, err
	}

	memories := append(*lastFewMemories, *aFewOlderMemories...)
	slices.Reverse(memories)

	return memories, nil
}

func NewMessage(incomingMessage *models.Message) *models.Message {

	return &models.Message{
		FromUserID:      models.GetSystemUser().ID,
		ToUserID:        incomingMessage.FromUserID,
		ReferenceID:     incomingMessage.ReferenceID,
		MessageType:     models.NewMessageTypeSMS(),
		MessageStatusID: models.NewMessageStatusSent().ID,
	}
}

func ValidateEvent(sqsEvent *events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {

		return fmt.Errorf("no event records found")
	}

	if sqsEvent.Records[0].Body == "" {

		return fmt.Errorf("no event body found")
	}

	return nil
}

func main() {

	log.New("SMS Sender Lambda booting.....").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config")
	}

	database := db.Get(cfg)

	if err := utils.PingDatabase(database); err != nil {
		log.New("Error pinging database").AddError(err).Log()

		return
	}

	if err := utils.PingGoogle(); err != nil {
		log.New("Error pinging Google. Possible bad internet connection.").AddError(err).Log()

		return
	}

	memoryService := message.NewMemoryService(
		message.NewMessageRepository(database), maxMemories,
	)

	handler := &SendSMSLambdaHandler{
		CompletionService: &ai.OpenAICompletionService{},
		MemoryService:     memoryService,
	}

	handler.Init(database)

	log.New("SMS Sender Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
