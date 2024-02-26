package models

type UserResponse struct {
	ID            int64  `json:"id"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Status        string `json:"status"`
	StatusID      int64  `json:"status_id"`
	PhoneVerified bool   `json:"phone_verified"`
}

func MakeUserResponseFromUser(user *User) *UserResponse {

	return &UserResponse{
		ID:            user.ID,
		Firstname:     user.Firstname,
		Lastname:      user.Lastname,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Status:        user.AccountStatus.Name,
		StatusID:      user.AccountStatusID,
		PhoneVerified: user.PhoneVerified,
	}
}
