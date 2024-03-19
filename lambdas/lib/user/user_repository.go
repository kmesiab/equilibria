package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/lib/status"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// UserRepository is a repository for managing Users.
type UserRepository struct {
	db               *gorm.DB
	statusRepository *status.Repository
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {

	return &UserRepository{
		db:               db,
		statusRepository: status.NewStatusRepository(db),
	}
}

// Create creates a new user in the database.
func (repo *UserRepository) Create(user *models.User) error {

	err := repo.db.Create(user).Error

	if err != nil {
		return err
	}

	statusPending, err := repo.statusRepository.FindByName("Pending Activation")

	if err != nil {
		return err
	}

	user.AccountStatus = *statusPending
	user.AccountStatusID = statusPending.ID

	return nil
}

// Update updates a user's details in the database.
func (repo *UserRepository) Update(user *models.User) error {

	return repo.db.Omit("AccountStatus").Updates(user).Error
}

// Delete deletes a user from the database.
func (repo *UserRepository) Delete(id int64) error {
	return repo.db.Delete(&models.User{}, id).Error
}

// ListAll lists all users in the database.
func (repo *UserRepository) ListAll() ([]models.User, error) {
	var users []models.User
	err := repo.db.Find(&users).Error

	return users, err
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Preload("AccountStatus").
		Where("email = ?", email).
		First(&user).Error

	return &user, err
}

// FindByID finds a user by their ID.
func (repo *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := repo.db.Preload("AccountStatus").First(&user, id).Error

	return &user, err
}

// FindByPhoneNumber finds a user by their phone number.
func (repo *UserRepository) FindByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := repo.db.Preload("AccountStatus").
		Where("phone_number = ?", phoneNumber).
		First(&user).Error

	return &user, err
}

func (repo *UserRepository) FindByName(name string) (*models.User, error) {
	var user models.User
	err := repo.db.Preload("AccountStatus").
		Where("firstname LIKE ? OR lastname LIKE ?", "%"+name+"%", "%"+name+"%").
		First(&user).Error

	return &user, err
}

func (repo *UserRepository) CheckPassword(phone, password string) (bool, error) {

	user, err := repo.FindByPhoneNumber(phone)

	if err != nil || user == nil {
		return false, err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		// Password does not match
		return false, nil
	}

	// Password matches
	return true, nil
}

func (repo *UserRepository) GetUsersWithoutConversationsSince(timeLimit time.Time) (*[]models.User, error) {

	var users []models.User

	err := repo.db.Model(users).Where(`
	account_status_id = ? 
	AND NOT EXISTS 
		(
			SELECT 1 FROM messages 
			WHERE (messages.from_user_id = users.id or messages.to_user_id = users.id)
			AND messages.created_at > ?
			AND phone_verified = true
	)`,
		models.AccountStatusActive,
		timeLimit.UTC(),
	).Find(&users).
		Error

	if err != nil {
		return nil, err
	}

	return &users, nil
}
