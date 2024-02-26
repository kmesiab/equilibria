package status

import (
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

type Repository struct {
	DB *gorm.DB
}

func NewStatusRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// FindByID retrieves an AccountStatus by its ID.
func (repo *Repository) FindByID(id int) (*models.AccountStatus, error) {
	var accountStatus models.AccountStatus
	result := repo.DB.First(&accountStatus, id)
	return &accountStatus, result.Error
}

// FindByName retrieves an AccountStatus by its name.
func (repo *Repository) FindByName(name string) (*models.AccountStatus, error) {
	var accountStatus models.AccountStatus
	result := repo.DB.Where("name = ?", name).First(&accountStatus)
	return &accountStatus, result.Error
}

// GetAll retrieves all AccountStatus entries.
func (repo *Repository) GetAll() ([]models.AccountStatus, error) {
	var accountStatuses []models.AccountStatus
	result := repo.DB.Find(&accountStatuses)
	return accountStatuses, result.Error
}
