package ai

import "github.com/kmesiab/equilibria/lambdas/models"

type CompletionServiceInterface interface {
	GetCompletion(message, prompt string, memories *[]models.Message) (string, error)
	CleanCompletionText(completion string) string
}

type MockCompletionService struct{}

func (m *MockCompletionService) GetCompletion(_, _ string, _ *[]models.Message) (string, error) {
	return "dummy completion", nil
}
