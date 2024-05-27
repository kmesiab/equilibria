package ai

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/forPelevin/gomoji"
	"github.com/sashabaranov/go-openai"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/encoding"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type OpenAICompletionService struct {
	RemoveEmojis bool
}

func (o *OpenAICompletionService) CleanCompletionText(completion string) string {

	if !encoding.IsGSMEncoded(completion) {

		log.New("Detected non GSM encoded completion: %s", completion).Log()
	}

	completion = strings.Replace(completion, "’", "'", -1)
	completion = strings.Replace(completion, "—", "-", -1)
	completion = strings.Replace(completion, "! ?", "!", -1)

	if o.RemoveEmojis {
		completion = gomoji.RemoveEmojis(completion)

		log.New("Removed emojis from completion").Log()
	}

	completion = strings.TrimSpace(completion)

	return completion
}

func (o *OpenAICompletionService) GetEmbeddings(text string) ([]float32, error) {

	client := openai.NewClient(config.Get().OpenAIAPIKey)

	embeddingsReq := openai.EmbeddingRequest{
		Model: EmbeddingServiceModel,
		Input: text,
	}

	embeddingsResp, err := client.CreateEmbeddings(context.Background(), embeddingsReq)

	if err != nil {

		return nil, err
	}

	if len(embeddingsResp.Data) == 0 {

		return nil, fmt.Errorf("embeddings data slice was empty")
	}

	// Log the retrieved embeddings for audit
	log.New("Successfully retrieved embeddings").
		Add("embeddings_count", strconv.Itoa(len(embeddingsResp.Data[0].Embedding))).
		Add("embeddings_object", embeddingsResp.Data[0].Object).
		Log()

	return embeddingsResp.Data[0].Embedding, nil
}

func (o *OpenAICompletionService) GetCompletion(
	message, prompt string, memories *[]models.Message,
) (string, error) {

	// Metrics
	var (
		promptTokenCount,
		historyTokenCount,
		totalTokenCount int

		memoryDump = ""
	)

	var totalMemories = 0
	var messages []openai.ChatCompletionMessage

	if memories != nil && len(*memories) > 0 {
		totalMemories = len(*memories)

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
	}

	promptTokenCount = len(prompt)
	totalTokenCount = historyTokenCount + promptTokenCount

	// Add current prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
	})

	// Add the current message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	log.New("OpenAI Audit Trail: Sending prompt.").
		Add("prompt_char_count", strconv.Itoa(promptTokenCount)).
		Add("history_char_count", strconv.Itoa(historyTokenCount)).
		Add("total_char_count", strconv.Itoa(totalTokenCount)).
		Add("total_memories", strconv.Itoa(totalMemories)).
		Add("prompt", prompt).
		Log()

	client := openai.NewClient(config.Get().OpenAIAPIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:            config.Get().ChatModelName,
			Messages:         messages,
			Temperature:      config.Get().ChatModelTemperature,
			MaxTokens:        config.Get().ChatModelMaxCompletionTokens,
			FrequencyPenalty: config.Get().ChatModelFrequencyPenalty,
		},
	)

	if err != nil {
		return "", err
	}

	log.New("OpenAI Audit Trail: Received response.").
		Add("model", resp.Model).
		Add("completion_tokens", strconv.Itoa(resp.Usage.CompletionTokens)).
		Add("prompt_tokens", strconv.Itoa(resp.Usage.PromptTokens)).
		Add("total_tokens", strconv.Itoa(resp.Usage.TotalTokens)).
		Add("response_content", resp.Choices[0].Message.Content).
		Add("prompt", prompt).
		Log()

	return resp.Choices[0].Message.Content, nil
}
