package user

import (
	"time"

	"github.com/kmesiab/equilibria/lambdas/models"
)

// UserService provides services related to users.
type UserService struct {
	repo *UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo *UserRepository) *UserService {

	return &UserService{repo: repo}
}

// Create creates a new user.
func (service *UserService) Create(user *models.User) error {

	return service.repo.Create(user)
}

// GetUserByID retrieves a user by their ID.
func (service *UserService) GetUserByID(id int64) (*models.User, error) {

	return service.repo.FindByID(id)
}

func (service *UserService) GetUsersByProviderCode(code string) (*[]models.User, error) {
	return service.repo.FindByProviderCode(code)
}

// GetUserByPhoneNumber retrieves a user by their phone number.
func (service *UserService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {

	return service.repo.FindByPhoneNumber(phoneNumber)
}

// Update updates a user's details.
func (service *UserService) Update(user *models.User) error {

	return service.repo.Update(user)
}

// DeleteUser deletes a user.
func (service *UserService) DeleteUser(id int64) error {

	return service.repo.Delete(id)
}

// ListUsers lists all users.
func (service *UserService) ListUsers() ([]models.User, error) {

	return service.repo.ListAll()
}

func (service *UserService) SystemUser() (*models.User, error) {

	return service.GetUserByID(1)
}

func (service *UserService) GetUsersWithoutConversationsSince(since time.Time) (*[]models.User, error) {

	return service.repo.GetUsersWithoutConversationsSince(since)
}
