package models

// MessageType represents the message_types table in the database with GORM annotations.
type MessageType struct {
	ID                int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string  `gorm:"size:100;unique;not null" json:"name"`
	BillRateInCredits float64 `gorm:"not null;default:0" json:"bill_rate_in_credits"`
}

func NewMessageTypeSMS() MessageType {
	return MessageType{
		ID:                2,
		Name:              "SMS",
		BillRateInCredits: .5,
	}
}
