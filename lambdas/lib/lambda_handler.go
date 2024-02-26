package lib

import (
	"github.com/aws/aws-lambda-go/events"
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/conversation"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/user"
)

type LambdaHandler struct {
	DB                  *gorm.DB
	UserService         *user.UserService
	MessageService      *message.MessageService
	ConversationService *conversation.ConversationService
}

func (h *LambdaHandler) Init(db *gorm.DB) *LambdaHandler {

	h.DB = db
	h.UserService = user.NewUserService(user.NewUserRepository(db))
	h.MessageService = message.NewMessageService(message.NewMessageRepository(db))
	h.ConversationService = conversation.NewConversationService(conversation.NewConversationRepository(db))

	return h

}

func RespondWithError(msg string, err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers:    config.DefaultHttpHeaders,
		StatusCode: statusCode,
		Body:       log.New(msg).AddError(err).Write(),
	}, nil
}
