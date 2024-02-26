package models

import "gorm.io/gorm"

type User struct {
	AccountStatus   AccountStatus `json:"status" gorm:"foreignKey:AccountStatusID;references:ID"`
	ID              int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Password        *string       `gorm:"type:varchar(1024);not null" json:"password"`
	PhoneNumber     string        `gorm:"type:varchar(20);unique;not null" json:"phone_number"`
	PhoneVerified   bool          `gorm:"not null;default:false" json:"phone_verified"`
	Firstname       string        `gorm:"type:varchar(100)" json:"firstname"`
	Lastname        string        `gorm:"type:varchar(100)" json:"lastname"`
	Email           string        `gorm:"type:varchar(100)" json:"email"`
	AccountStatusID int64         `gorm:"not null;" json:"account_status_id"`
}

func (u *User) IsValid() bool {

	return u.PhoneNumber != "" &&
		u.Email != "" &&
		u.Firstname != "" &&
		u.Lastname != ""
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	// omit the data field if it is empty
	if u.AccountStatusID == 0 {
		tx.Statement.Omit("account_status_id")
	}

	if u.PhoneNumber == "" {
		tx.Statement.Omit("phone_number")
	}

	return
}

func GetSystemUser() *User {
	return &User{
		ID:              1,
		Firstname:       "System",
		Lastname:        "User",
		Email:           "-@-",
		PhoneNumber:     "+18333595081",
		AccountStatusID: 1,
		PhoneVerified:   true,
	}
}
