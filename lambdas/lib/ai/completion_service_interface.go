package ai

import "github.com/kmesiab/equilibria/lambdas/models"

type CompletionServiceInterface interface {
	GetCompletion(message, prompt string, memories *[]models.Message) (string, error)
	CleanCompletionText(completion string) string
	GetEmbeddings(text string) ([]float32, error)
}
