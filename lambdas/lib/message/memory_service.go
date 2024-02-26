package message

import (
	"fmt"

	"github.com/kmesiab/equilibria/lambdas/models"
)

type MemoryService struct {
	*MessageService
}

func NewMemoryService(repo *Repository) *MemoryService {

	svc := NewMessageService(repo)
	return &MemoryService{svc}
}

const (
	maxMemories   = 150
	promptHeader  = ""
	promptLineFmt = "[%s] %s: %s\n"
)

// GeneratePrompt generates a prompt from the user's previous messages.
// The prompt output is in the format specified by promptLnFmt. In general,
// the output appears like a chat log: "[date] [name]: [message]\n\n"
func (m *MemoryService) GeneratePrompt(user *models.User) (string, error) {

	var prompt = promptHeader
	messages, err := m.FindByUser(user)

	if err != nil {
		return "", err
	}

	var recentMemories []models.Message

	// Take only the last maxMemories messages
	if len(*messages) > maxMemories {
		recentMemories = (*messages)[len(*messages)-maxMemories:]
	} else {
		recentMemories = *messages
	}

	for _, message := range recentMemories {
		prompt += fmt.Sprintf(promptLineFmt,
			message.CreatedAt,
			message.From.Firstname,
			message.Body,
		)
	}

	return prompt, nil
}

// GetMemories returns a slice of Messages for this user, representing all
// memories for all conversations between this user and the service.
func (m *MemoryService) GetMemories(user *models.User) (*[]models.Message, error) {

	return m.GetRandomMessagePairs(user, maxMemories)
}
