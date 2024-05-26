package ai

import "github.com/kmesiab/equilibria/lambdas/models"

type MockCompletionService struct{}

func (m *MockCompletionService) GetCompletion(_, _ string, _ *[]models.Message) (string, error) {
	return "dummy completion", nil
}

func (m *MockCompletionService) CleanCompletionText(_ string) string {
	return "dummy cleaned text"
}

func (m *MockCompletionService) GetEmbeddings(_ string) ([]float32, error) {
	return []float32{0.0, 1.0, 2.0}, nil
}
