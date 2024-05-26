package fact_agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
)

func TestFactAgent_Do(t *testing.T) {

	completionSvc := &ai.OpenAICompletionService{RemoveEmojis: true}
	agent := NewFactAgent(completionSvc)

	input := "My dog died today."

	output, err := agent.Do(input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}
