package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TransactionTypeStringCredit = "credit"
	TransactionTypeStringDebit  = "debit"
	FundingSourceStringStripe   = "stripe"
)

// Transaction represents a credit transaction in the database.
type Transaction struct {
	gorm.Model
	User            User         `gorm:"foreignKey:UserID"`
	Conversation    Conversation `gorm:"foreignKey:ConversationID"`
	ID              int64        `gorm:"primaryKey; autoIncrement"`
	UserID          int64        `gorm:"not null; index"`
	ConversationID  int64        `gorm:"not null; index"`
	Amount          float64      `gorm:"type:decimal(10,2); not null"`
	TransactionType string       `gorm:"type:enum('credit', 'debit'); not null"`
	FundingSource   string       `gorm:"type:enum('stripe', 'paypal', 'bank_transfer', 'cash', 'refund', 'customer credit'); not null"`
	Description     string       `gorm:"type:text"`
	ReferenceID     string       `gorm:"type:varchar(255)"`
	Timestamp       *time.Time   `gorm:"type:datetime; default:CURRENT_TIMESTAMP"`
}
