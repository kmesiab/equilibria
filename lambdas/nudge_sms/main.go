package main

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/kmesiab/equilibria/lambdas/lib"
	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/twilio"
	"github.com/kmesiab/equilibria/lambdas/lib/user"
	"github.com/kmesiab/equilibria/lambdas/lib/utils"
	"github.com/kmesiab/equilibria/lambdas/models"
)

const (
	HoursSinceLastNudge                   = 7
	MaxNewMemories                        = 20
	MaxOldMemories                        = 15
	NumMemoriesToBeConsideredExistingUser = 3
)

var TimeSinceLastMessage time.Time

type NudgeSMSLambdaHandler struct {
	lib.LambdaHandler

	UserService       *user.UserService
	MemoryService     *message.MemoryService
	CompletionService ai.CompletionServiceInterface

	MaxNewMemories         int
	MaxOldMemories         int
	NudgeIfNoMessagesSince time.Time
}

func (h *NudgeSMSLambdaHandler) HandleRequest(e events.EventBridgeEvent) error {

	log.New(
		"Looking up users without conversations since %s",
		h.NudgeIfNoMessagesSince.Format("2006-01-02 15:04:05"),
	).Log()

	users, err := h.UserService.GetUsersWithoutConversationsSince(h.NudgeIfNoMessagesSince)

	if err != nil {
		log.New("Error looking up users without conversations").AddError(err).Log()

		return err
	}

	if len(*users) == 0 {
		log.New("No users without conversations found.  Exiting.").Log()

		return nil
	}

	wg := &sync.WaitGroup{}

	for _, u := range *users {

		if !u.NudgesEnabled() {

			log.New("User %s is not nudging", u.PhoneNumber).Log()

			continue
		}

		wg.Add(1)

		// Nudge each user in a goroutine
		go func(u *models.User, wg *sync.WaitGroup) {
			defer wg.Done()

			// Skip the system user
			if u.PhoneNumber == models.GetSystemUser().PhoneNumber {

				return
			}

			err := h.Nudge(u)

			// Nudge failed for some reason
			if err != nil {

				log.New("Error nudging user %s", u.PhoneNumber).Log()
			}
		}(&u, wg)
	}

	log.New("Waiting for %d nudge SMS routines to complete", len(*users)).Log()

	wg.Wait()

	log.New("All Nudges completed. %d users nudged.", len(*users)).Log()

	return nil
}

func (h *NudgeSMSLambdaHandler) Nudge(user *models.User) error {

	memories, err := h.GetMemories(user)

	slices.Reverse(*memories)

	if err != nil {
		log.New("Error retrieving memories for user %s", user.PhoneNumber).
			AddError(err).AddUser(user).Log()

		return err
	}

	memoryDumpString := MemoriesToString(memories)

	// The number of messages I've sent to the system.
	myMemories := utils.FilterSlice(*memories, func(m models.Message) bool {
		return m.FromUserID != models.GetSystemUser().ID
	})

	var promptModifier string

	log.New("Total User Texts: %d", len(myMemories)).Log()

	if len(myMemories) < NumMemoriesToBeConsideredExistingUser {
		log.New("Using new user prompt modifier").AddUser(user).Log()

		promptModifier = NudgePromptNewUserModifier
	} else {
		log.New("Using existing user prompt modifier").AddUser(user).Log()

		promptModifier = NudgePromptExistingUserModifier
	}

	prompt := fmt.Sprintf(NudgePrompt, promptModifier, user.Firstname)

	log.New("Attaching %d memories", len(*memories)).
		Add("memory_dump", memoryDumpString).
		Add("prompt", prompt).
		AddUser(user).
		Log()

	completion, err := h.CompletionService.GetCompletion(prompt, prompt, memories)

	if err != nil {
		log.New("Error retrieving a few older memories for user %s", user.PhoneNumber).
			AddError(err).AddUser(user).Log()

		return err
	}

	var (
		now        = time.Now()
		convo      *models.Conversation
		newMessage *models.Message
	)

	log.New("Starting conversation for %s", user.PhoneNumber).AddUser(user).Log()

	// Start the conversation
	if convo, err = h.CreateConversation(user, completion, now); err != nil {
		log.New("Error closing nudge conversation for %s", user.PhoneNumber).
			Add("completion", completion).
			AddUser(user).
			AddError(err).
			Log()

		return err
	}

	log.New("Sending nudge SMS to %s", user.PhoneNumber).
		Add("conversation_id", strconv.FormatInt(convo.ID, 10)).
		Add("completion", completion).
		AddUser(user).
		Log()

	smsResponse, err := h.SendSms(user, completion)

	if err != nil {
		log.New("Error sending SMS for %s: %s", user.PhoneNumber, err.Error()).
			Add("completion", completion).
			AddError(err).
			AddUser(user).
			Log()

		return err
	}

	// Add this message to the conversation
	if newMessage, err = h.CreateMessage(user, convo, completion, *smsResponse.Sid); err != nil {
		log.New("Error creating message for %s", user.PhoneNumber).
			Add("completion", completion).
			AddUser(user).
			AddError(err).
			Log()

		return err
	}

	log.New("Closing conversation for %s", user.PhoneNumber).
		Add("conversation_id", strconv.FormatInt(convo.ID, 10)).
		Add("completion", completion).
		AddMessage(newMessage).
		AddUser(user).
		Log()

	err = h.ConversationService.EndConversation(convo.ID)

	if err != nil {
		log.New("Error closing nudge conversation for %s", user.PhoneNumber).
			Add("completion", completion).
			AddMessage(newMessage).
			AddUser(user).
			Log()

		return err
	}

	return nil
}

func (h *NudgeSMSLambdaHandler) GetMemories(user *models.User) (*[]models.Message, error) {

	lastFewMemories, err := h.MemoryService.GetLastNMessagePairs(user, h.MaxNewMemories)

	if err != nil {
		log.New(
			"Error retrieving last few memories for user %s",
			user.PhoneNumber).Log()

		return nil, err
	}

	aFewOlderMemories, err := h.MemoryService.GetRandomMessagePairs(user, h.MaxOldMemories)

	if err != nil {
		log.New(
			"Error retrieving a few older memories for user %s",
			user.PhoneNumber).Log()

		return nil, err
	}

	memories := append(*lastFewMemories, *aFewOlderMemories...)

	return &memories, nil
}

func (h *NudgeSMSLambdaHandler) SendSms(recipient *models.User, completion string) (*openapi.ApiV2010Message, error) {

	smsResponse, err := twilio.SendSMS(
		models.GetSystemUser().PhoneNumber, recipient.PhoneNumber, completion,
	)

	if err != nil {

		return nil, err
	}

	if smsResponse == nil {

		return nil, fmt.Errorf("error SMS response is empty")
	}

	log.New("Successfully sent outbound nudge message from %s to %s",
		models.GetSystemUser().PhoneNumber, recipient.PhoneNumber,
	).
		Add("status", *smsResponse.Status).
		Add("sms_message_id", *smsResponse.Sid).
		Add("num_segments", *smsResponse.NumSegments).
		Log()

	return smsResponse, nil
}

func (h *NudgeSMSLambdaHandler) CreateMessage(recipient *models.User, conversation *models.Conversation, completion, refID string) (*models.Message, error) {

	var now = time.Now()

	newMessage := &models.Message{
		FromUserID:      models.GetSystemUser().ID,
		ToUserID:        recipient.ID,
		ReferenceID:     &refID,
		MessageType:     models.NewMessageTypeSMS(),
		MessageStatusID: models.NewMessageStatusSent().ID,
		MessageStatus:   models.NewMessageStatusSending(),
		Body:            completion,
		SentAt:          &now,
		ConversationID:  conversation.ID,
		To:              *recipient,
	}

	if err := h.MessageService.CreateMessage(newMessage); err != nil {

		return nil, err
	}

	return newMessage, nil
}

func (h *NudgeSMSLambdaHandler) CreateConversation(recipient *models.User, completion string, now time.Time) (*models.Conversation, error) {

	// Start a new conversation
	convo := &models.Conversation{
		User:      *models.GetSystemUser(),
		UserID:    models.GetSystemUser().ID,
		StartTime: &now,
	}

	err := h.ConversationService.CreateConversation(convo)

	if err != nil || convo.ID == 0 {
		log.New("Error creating a new conversation for %s", recipient.PhoneNumber).
			Add("completion", completion).
			AddUser(recipient).
			AddError(err).
			Log()

		return nil, err
	}
	return convo, nil
}

func main() {

	log.New("SMS Nudger Lambda booting.....").Log()

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
		log.New("Error pinging Google. Possible bad internet connection.").
			AddError(err).Log()

		return
	}

	TimeSinceLastMessage = time.Now().UTC().Add(time.Hour * -HoursSinceLastNudge)

	usrSvc := user.NewUserService(
		user.NewUserRepository(database),
	)

	memSvc := message.NewMemoryService(
		message.NewMessageRepository(database), MaxNewMemories+MaxOldMemories,
	)

	llmSvc := &ai.OpenAICompletionService{
		RemoveEmojis: false,
	}

	handler := &NudgeSMSLambdaHandler{
		UserService:            usrSvc,
		MemoryService:          memSvc,
		CompletionService:      llmSvc,
		MaxNewMemories:         MaxNewMemories,
		MaxOldMemories:         MaxOldMemories,
		NudgeIfNoMessagesSince: TimeSinceLastMessage,
	}

	handler.Init(database)

	log.New("SMS Nudger Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}

func MemoriesToString(memories *[]models.Message) string {

	var memoriesString = ""

	for _, m := range *memories {
		memoriesString += fmt.Sprintf("[%s] %s: %s\n", m.CreatedAt, m.From.Firstname, m.Body)
	}

	return memoriesString
}
