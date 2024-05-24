package message

import (
	"fmt"

	"github.com/kmesiab/equilibria/lambdas/models"
)

const (
	promptHeader  = ""
	promptLineFmt = "[%s] %s: %s\n"
)

type MemoryService struct {
	*MessageService
}

func NewMemoryService(repo *Repository) *MemoryService {

	svc := NewMessageService(repo)
	return &MemoryService{MessageService: svc}
}

// GeneratePrompt generates a prompt from the user's previous messages.
// The prompt output is in the format specified by promptLnFmt. In general,
// the output appears like a chat log: "[date] [name]: [message]\n\n"
func (m *MemoryService) GeneratePrompt(user *models.User) (string, error) {

	var prompt = promptHeader
	messages, err := m.FindByUser(user)

	if err != nil {
		return "", err
	}

	for _, message := range *messages {
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
func (m *MemoryService) GetMemories(user *models.User, limit int) (*[]models.Message, error) {

	return m.GetRandomMessagePairs(user, limit)
}

func (m *MemoryService) GetLastNMessagePairs(user *models.User, size int) (*[]models.Message, error) {

	return m.repo.GetLastNMessagePairs(user, size)
}
