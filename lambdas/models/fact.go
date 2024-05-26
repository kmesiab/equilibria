package models

import (
	"time"
)

// Fact represents a fact in the database.
type Fact struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64      `json:"user_id" gorm:"not null;index"`
	ConversationID int64      `json:"conversation_id" gorm:"not null;index"`
	Body           string     `json:"body" gorm:"type:text"`
	Reasoning      string     `json:"reasoning" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"type:datetime;default:null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
}
