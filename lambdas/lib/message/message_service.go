package message

import (
	"github.com/kmesiab/equilibria/lambdas/models"
)

// MessageService provides services related to messages.
type MessageService struct {
	repo *Repository
}

// NewMessageService New creates a new MessageService.
func NewMessageService(repo *Repository) *MessageService {

	return &MessageService{repo: repo}
}

// CreateMessage creates a new message.
func (service *MessageService) CreateMessage(message *models.Message) error {

	return service.repo.Create(message)
}

// FindByID retrieves a message by its ID.
func (service *MessageService) FindByID(id int64) (*models.Message, error) {

	return service.repo.FindByID(id)
}

// UpdateMessage updates a message's details.
func (service *MessageService) UpdateMessage(message *models.Message) error {

	return service.repo.Update(message)
}

// UpdateStatus updates a message's status.
func (service *MessageService) UpdateStatus(message *models.Message) error {

	return service.repo.UpdateStatus(message)
}

// DeleteMessage deletes a message.
func (service *MessageService) DeleteMessage(id int64) error {

	return service.repo.Delete(id)
}

func (service *MessageService) GetRandomMessagePairs(user *models.User, limit int) (*[]models.Message, error) {

	return service.repo.GetRandomMessagePairs(user, limit)
}

func (service *MessageService) FindByUser(user *models.User) (*[]models.Message, error) {

	return service.repo.FindByUser(user)
}

func (service *MessageService) FindByReferenceID(refID string) (*models.Message, error) {

	return service.repo.FindByReferenceID(refID)
}

func (service *MessageService) FindConversationIDByReferenceID(refID string) (int64, error) {

	msg, err := service.FindByReferenceID(refID)

	if err != nil {
		return 0, err
	}

	return msg.ConversationID, err

}
