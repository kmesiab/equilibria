package models

import (
	"time"
)

type NrcLex struct {
	ID            int64      `gorm:"primaryKey" json:"id"`
	UserID        int64      `json:"user_id"`
	MessageID     int64      `json:"message_id"`
	Anger         float64    `json:"anger"`
	Anticipation  float64    `json:"anticipation"`
	Disgust       float64    `json:"disgust"`
	Fear          float64    `json:"fear"`
	Trust         float64    `json:"trust"`
	Joy           float64    `json:"joy"`
	Negative      float64    `json:"negative"`
	Positive      float64    `json:"positive"`
	Sadness       float64    `json:"sadness"`
	Surprise      float64    `json:"surprise"`
	VaderCompound float64    `json:"vader_compound"`
	VaderNeg      float64    `json:"vader_neg"`
	VaderNeu      float64    `json:"vader_neu"`
	VaderPos      float64    `json:"vader_pos"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (NrcLex) TableName() string {
	return "nrclex"
}
