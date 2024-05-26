package facts

import (
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

// Repository is a repository for managing Facts.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// SetDBInstance sets the database instance for the repository.
func (r *Repository) SetDBInstance(db *gorm.DB) {
	r.db = db
}

// FindByID retrieves a fact by its ID from the database.
func (r *Repository) FindByID(id int64) (*models.Fact, error) {
	var fact models.Fact
	if err := r.db.First(&fact, id).Error; err != nil {
		return nil, err
	}
	return &fact, nil
}

// Create inserts a new fact into the database.
func (r *Repository) Create(fact *models.Fact) error {
	return r.db.Create(fact).Error
}

// Update modifies an existing fact in the database.
func (r *Repository) Update(fact *models.Fact) error {
	return r.db.Save(fact).Error
}

// Delete performs a soft delete on a fact by its ID.
func (r *Repository) Delete(id int64) error {
	return r.db.Delete(&models.Fact{}, id).Error
}
