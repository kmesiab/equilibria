package conversation

import (
	"github.com/kmesiab/equilibria/lambdas/models"
)

// ConversationService provides services related to conversations.
type ConversationService struct {
	repo *ConversationRepository
}

// NewConversationService creates a new ConversationService.
func NewConversationService(repo *ConversationRepository) *ConversationService {

	return &ConversationService{repo: repo}
}

// StartConversation creates and starts a new conversation.
func (service *ConversationService) StartConversation(conversationID int64) error {

	return service.repo.StartConversation(conversationID)
}

// EndConversation ends an existing conversation.
func (service *ConversationService) EndConversation(conversationID int64) error {

	return service.repo.EndConversation(conversationID)
}

// FindByID retrieves a conversation by its ID.
func (service *ConversationService) FindByID(id int64) (*models.Conversation, error) {

	return service.repo.FindByID(id)
}

// FindByUser retrieves a conversation by the user.
func (service *ConversationService) FindByUser(user models.User) (*[]models.Conversation, error) {

	return service.repo.FindByUser(user)
}

// GetAllConversations retrieves all conversations.
func (service *ConversationService) GetAllConversations() ([]models.Conversation, error) {

	return service.repo.GetAll()
}

// UpdateConversation updates the details of an existing conversation.
func (service *ConversationService) UpdateConversation(conversation *models.Conversation) error {

	return service.repo.Update(conversation)
}

// CreateConversation adds a new conversation to the database.
func (service *ConversationService) CreateConversation(conversation *models.Conversation) error {

	return service.repo.Create(conversation)
}

// SoftDeleteConversation marks a conversation as deleted.
func (service *ConversationService) SoftDeleteConversation(id int64) error {

	return service.repo.SoftDelete(id)
}

// HardDeleteConversation permanently deletes a conversation.
func (service *ConversationService) HardDeleteConversation(id int64) error {

	return service.repo.HardDelete(id)
}

func (service *ConversationService) GetOpenConversationsByUserID(userID int64) (*[]models.Conversation, error) {

	return service.repo.GetOpenConversationsByUserID(userID)
}
