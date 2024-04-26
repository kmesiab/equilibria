package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanCompletionText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Replace apostrophes",
			input:    "It’s a test",
			expected: "It's a test",
		},
		{
			name:     "Replace '! ?' with '!'",
			input:    "Really! ?",
			expected: "Really!",
		},
		{
			name:     "Remove emojis when flag is true",
			input:    "Hello 👋 world",
			expected: "Hello  world",
		},
		{
			name:     "Trim whitespace",
			input:    "  trimmed  ",
			expected: "trimmed",
		},
	}

	o := &OpenAICompletionService{
		RemoveEmojis: true,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := o.CleanCompletionText(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCleanCompletionText_NoEmojiRemoval(t *testing.T) {
	// Assuming RemoveEmojis is a property of OpenAICompletionService in this context
	o := &OpenAICompletionService{RemoveEmojis: false}

	input := "Hello 👋 world"
	expected := "Hello 👋 world" // Expecting the input to remain unchanged regarding emojis

	result := o.CleanCompletionText(input)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
