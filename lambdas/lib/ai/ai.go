package ai

import "github.com/kmesiab/equilibria/lambdas/models"

type CompletionServiceInterface interface {
	GetCompletion(prompt string, memories *[]models.Message) (string, error)
}

type MockCompletionService struct{}

func (m *MockCompletionService) GetCompletion(_ string, _ *[]models.Message) (string, error) {
	return "dummy completion", nil
}
