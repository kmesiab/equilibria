package models

type UserResponse struct {
	ID            int64  `json:"id"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Status        string `json:"status"`
	StatusID      int64  `json:"status_id"`
	UserTypeID    int64  `json:"user_type_id"`
	PhoneVerified bool   `json:"phone_verified"`
	NudgeEnabled  bool   `json:"nudge_enabled"`
	ProviderCode  string `json:"provider_code"`
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
		UserTypeID:    user.UserTypeID,
		PhoneVerified: user.PhoneVerified,
		NudgeEnabled:  user.NudgesEnabled(),
		ProviderCode:  user.ProviderCode,
	}
}
