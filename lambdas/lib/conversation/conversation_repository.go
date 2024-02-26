package conversation

import (
	"time"

	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

type ConversationRepository struct {
	DB *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{DB: db}
}

// StartConversation creates a new Conversation record.
func (repo *ConversationRepository) StartConversation(conversationID int64) error {
	now := time.Now()
	conversation := &models.Conversation{
		ID:        conversationID,
		StartTime: &now,
	}

	return repo.Update(conversation)
}

// EndConversation sets the conversation as ended.
func (repo *ConversationRepository) EndConversation(conversationID int64) error {
	now := time.Now()
	conversation := &models.Conversation{
		ID:      conversationID,
		EndTime: &now,
	}

	return repo.Update(conversation)
}

// FindByID retrieves a Conversation by its ID.
func (repo *ConversationRepository) FindByID(conversationID int64) (*models.Conversation, error) {
	var conversation models.Conversation
	// TODO: Consider batching
	result := repo.DB.
		Preload("User").
		Preload("User.AccountStatus").
		First(&conversation, conversationID)

	return &conversation, result.Error
}

// FindByUser retrieves a Conversation by its user ID.
func (repo *ConversationRepository) FindByUser(user models.User) (*[]models.Conversation, error) {
	var conversation []models.Conversation

	result := repo.DB.Preload("User.AccountStatus").
		Where("user_id = ?", user.ID).Find(&conversation)

	return &conversation, result.Error
}

// GetAll retrieves all Conversation entries.
func (repo *ConversationRepository) GetAll() ([]models.Conversation, error) {
	var conversations []models.Conversation
	result := repo.DB.Preload("User.AccountStatus").Find(&conversations)

	return conversations, result.Error
}

// Update updates an existing Conversation.
func (repo *ConversationRepository) Update(conversation *models.Conversation) error {

	return repo.DB.Updates(conversation).Error
}

// Create adds a new Conversation to the database.
func (repo *ConversationRepository) Create(conversation *models.Conversation) error {

	return repo.DB.Create(conversation).Error
}

// SoftDelete marks a Conversation as deleted.
func (repo *ConversationRepository) SoftDelete(id int64) error {

	return repo.DB.Delete(&models.Conversation{}, id).Error
}

// HardDelete permanently deletes a Conversation record.
func (repo *ConversationRepository) HardDelete(id int64) error {

	return repo.DB.Unscoped().Where("id = ?", id).
		Delete(&models.Conversation{}).Error
}

func (repo *ConversationRepository) GetOpenConversationsByUserID(userID int64) (*[]models.Conversation, error) {
	var conversations []models.Conversation
	result := repo.DB.Preload("User.AccountStatus").
		Where("user_id = ? AND end_time IS NULL", userID).
		Find(&conversations)

	return &conversations, result.Error
}
