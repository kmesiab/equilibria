package models

import (
	"time"

	"gorm.io/gorm"
)

// Conversation represents the structure of the 'conversations' table.
// In practice, a conversation is a collection of one or two messages.
// If a Conversation has one message, then it should be safe to assume
// the EndTime is null. If a Conversation has two messages, then then
// the conversation is considered closed, and the EndTime should be set.
type Conversation struct {
	User      User       `gorm:"foreignKey:UserID" json:"user"`
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"notNull" json:"user_id"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`
}

func (c *Conversation) BeforeUpdate(tx *gorm.DB) (err error) {
	if c.StartTime == nil {
		tx.Statement.Omit("start_time")
	}

	return
}
