package message_type

import (
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

type MessageTypeRepository struct {
	DB *gorm.DB
}

func NewMessageTypeRepository(db *gorm.DB) *MessageTypeRepository {
	return &MessageTypeRepository{DB: db}
}

// FindByID retrieves a MessageType by its ID.
func (repo *MessageTypeRepository) FindByID(id int) (*models.MessageType, error) {
	var messageType models.MessageType
	result := repo.DB.First(&messageType, id)
	return &messageType, result.Error
}

// FindByName retrieves a MessageType by its name.
func (repo *MessageTypeRepository) FindByName(name string) (*models.MessageType, error) {
	var messageType models.MessageType
	result := repo.DB.Where("name = ?", name).First(&messageType)
	return &messageType, result.Error
}

// GetAll retrieves all MessageType entries.
func (repo *MessageTypeRepository) GetAll() ([]models.MessageType, error) {
	var messageTypes []models.MessageType
	result := repo.DB.Find(&messageTypes)
	return messageTypes, result.Error
}
