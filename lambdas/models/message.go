package models

import (
	"time"

	"gorm.io/gorm"
)

// Message represents the messages table in the database with GORM annotations.
type Message struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ReferenceID     *string        `gorm:"size:255" json:"reference_id"`
	ConversationID  int64          `gorm:"index:idx_conversation,sort:asc;foreignKey" json:"conversation_id"`
	FromUserID      int64          `gorm:"not null;foreignKey" json:"from_user_id"`
	ToUserID        int64          `gorm:"not null;foreignKey" json:"to_user_id"`
	Body            string         `gorm:"not null;size:255" json:"body"`
	MessageTypeID   int64          `gorm:"not null;foreignKey" json:"message_type_id"`
	MessageStatusID int64          `gorm:"not null;foreignKey" json:"message_status_id"`
	SentAt          *time.Time     `gorm:"default:null" json:"sent_at"`
	ReceivedAt      *time.Time     `gorm:"default:null" json:"received_at"`
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_reference_id,sort:asc;default:null" json:"deleted_at"`

	// Foreign key relationships
	Conversation  Conversation  `gorm:"foreignKey:ConversationID" json:"conversation"`
	MessageStatus MessageStatus `gorm:"foreignKey:MessageStatusID;association_autoupdate:false;association_autocreate:false" json:"message_status"`
	MessageType   MessageType   `gorm:"foreignKey:MessageTypeID;association_autoupdate:false;association_autocreate:false" json:"message_type"`
	From          User          `gorm:"foreignKey:FromUserID;association_autoupdate:false;association_autocreate:false" json:"from_user"`
	To            User          `gorm:"foreignKey:ToUserID;association_autoupdate:false;association_autocreate:false" json:"to_user"`
}
