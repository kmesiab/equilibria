package message

import (
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

// Repository provides an interface to perform CRUD operations on Message entities.
type Repository struct {
	DB *gorm.DB
}

// NewMessageRepository creates a new instance of MessageRepository.
func NewMessageRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// Create inserts a new Message into the database.
func (r *Repository) Create(message *models.Message) error {
	return r.DB.Create(message).Error
}

// FindByID finds a Message by its ID.
func (r *Repository) FindByID(id int64) (*models.Message, error) {
	var message models.Message
	err := r.DB.Preload("Conversation.User.AccountStatus").
		Preload("MessageType").
		Preload("MessageStatus").
		Preload("From.AccountStatus").
		Preload("To.AccountStatus").
		Preload("To.AccountStatus").
		Where("id = ?", id).
		First(&message).Error

	return &message, err
}

// GetRandomMessagePairs finds a Message by its ID.
func (r *Repository) GetRandomMessagePairs(user *models.User, limit int) (*[]models.Message, error) {
	var messages []models.Message

	err := r.DB.Raw(`
		SELECT * FROM messages 
		WHERE (from_user_id = ? or to_user_id = ?)
		AND messages.conversation_id IN (
			SELECT id FROM conversations ORDER BY RAND()
		)
		ORDER BY created_at, from_user_id DESC 
		LIMIT ?
	`, user.ID, user.ID, limit).Scan(&messages).Error

	if err != nil {
		return nil, err
	}

	return &messages, nil
}

// FindByUser finds a Message by its ID.
func (r *Repository) FindByUser(user *models.User) (*[]models.Message, error) {
	var messages []models.Message
	err := r.DB.Preload("Conversation.User.AccountStatus").
		Preload("MessageType").
		Preload("MessageStatus").
		Preload("From.AccountStatus").
		Preload("To.AccountStatus").
		Preload("To.AccountStatus").
		// Where("from_user_id = ? or to_user_id = ?", user.ID, user.ID).
		Where("from_user_id = ?", user.ID).
		Find(&messages).Error

	return &messages, err
}

func (r *Repository) FindByReferenceID(refID string) (*models.Message, error) {
	var message = &models.Message{}
	err := r.DB.Preload("Conversation.User.AccountStatus").
		Preload("MessageType").
		Preload("MessageStatus").
		Preload("From.AccountStatus").
		Preload("To.AccountStatus").
		Preload("To.AccountStatus").
		Where("reference_id = ?", refID).
		First(message).Error

	return message, err
}

// Update updates an existing Message.
func (r *Repository) Update(message *models.Message) error {

	return r.DB.Model(&message).Omit("id").Updates(message).Error
}

func (r *Repository) UpdateStatus(message *models.Message) error {

	return r.DB.Model(&message).
		Select("message_status_id").
		Updates(message).
		Error
}

// Delete removes a Message from the database.
func (r *Repository) Delete(id int64) error {
	return r.DB.Delete(&models.Message{}, id).Error
}
