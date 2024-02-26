package ai

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sashabaranov/go-openai"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/models"
)

const CompletionServiceModel = openai.GPT4TurboPreview

type OpenAICompletionService struct{}

func (o *OpenAICompletionService) GetCompletion(
	prompt string, memories *[]models.Message,
) (string, error) {

	// Metrics
	var (
		promptTokenCount,
		historyTokenCount,
		totalTokenCount,
		totalMemories int

		memoryDump = ""
	)

	var messages []openai.ChatCompletionMessage

	// Add memories
	for _, m := range *memories {

		role := openai.ChatMessageRoleUser

		if m.FromUserID == 1 {
			role = openai.ChatMessageRoleAssistant
		}

		body := fmt.Sprintf("%s %s", m.CreatedAt, m.Body)
		memoryDump += body + "\n"

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: body,
		})

		// Keep count of the tokens used in just the
		// memories, and the primary payload.
		historyTokenCount += len(m.Body)
	}

	totalMemories = len(*memories)
	promptTokenCount = len(prompt)
	totalTokenCount = historyTokenCount + promptTokenCount

	// Add current prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	log.New("OpenAI Audit Trail: Sending prompt.").
		Add("prompt_char_count", strconv.Itoa(promptTokenCount)).
		Add("history_char_count", strconv.Itoa(historyTokenCount)).
		Add("total_char_count", strconv.Itoa(totalTokenCount)).
		Add("total_memories", strconv.Itoa(totalMemories)).
		Add("prompt", prompt).
		Add("memory_dump", memoryDump).
		Log()

	client := openai.NewClient(config.Get().OpenAIAPIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    CompletionServiceModel,
			Messages: messages,
		},
	)

	log.New("OpenAI Audit Trail: Received response.").
		Add("model", resp.Model).
		Add("completion_tokens", strconv.Itoa(resp.Usage.CompletionTokens)).
		Add("prompt_tokens", strconv.Itoa(resp.Usage.PromptTokens)).
		Add("total_tokens", strconv.Itoa(resp.Usage.TotalTokens)).
		Add("response_content", resp.Choices[0].Message.Content).
		Add("prompt", prompt).
		Log()

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
