package nrclex

import (
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

// Repository provides an interface to perform CRUD operations on NrcLex entities.
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of NrcLexRepository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// Create inserts a new NrcLex into the database.
func (r *Repository) Create(nrcLex *models.NrcLex) error {
	return r.DB.Create(nrcLex).Error
}

// FindByID finds a NrcLex by its ID.
func (r *Repository) FindByID(id int64) (*models.NrcLex, error) {
	var nrcLex models.NrcLex

	err := r.DB.First(&nrcLex, id).Error

	if err != nil {
		return nil, err
	}

	return &nrcLex, nil
}

// Update updates an existing NrcLex.
func (r *Repository) Update(nrcLex *models.NrcLex) error {
	// This starts a new DB transaction.
	tx := r.DB.Begin()

	// Attempt to update the record that matches the given ID.
	// This will only update existing records and won't insert a new one if the record doesn't exist.
	result := tx.Model(&models.NrcLex{}).Where("id = ?", nrcLex.ID).Updates(nrcLex)

	// Check if the record was found and updated. If not, rollback and return an error.
	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// Check for other possible errors during the update.
	if err := result.Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything was okay.
	return tx.Commit().Error
}

// Delete removes a NrcLex from the database.
func (r *Repository) Delete(id int64) error {
	return r.DB.Delete(&models.NrcLex{}, id).Error
}
