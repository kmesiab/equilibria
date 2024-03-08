package nrclex

import "github.com/kmesiab/equilibria/lambdas/models"

// Service provides services related to NrcLex entities.
type Service struct {
	repo *Repository
}

// NewService creates a new Service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateNrcLex creates a new NrcLex.
func (service *Service) CreateNrcLex(nrcLex *models.NrcLex) error {
	return service.repo.Create(nrcLex)
}

// FindNrcLexByID retrieves a NrcLex by its ID.
func (service *Service) FindNrcLexByID(id int64) (*models.NrcLex, error) {
	return service.repo.FindByID(id)
}

// UpdateNrcLex updates an existing NrcLex.
func (service *Service) UpdateNrcLex(nrcLex *models.NrcLex) error {
	return service.repo.Update(nrcLex)
}

// DeleteNrcLex deletes a NrcLex.
func (service *Service) DeleteNrcLex(id int64) error {
	return service.repo.Delete(id)
}
