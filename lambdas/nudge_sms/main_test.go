package main_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/models"
	"github.com/kmesiab/equilibria/lambdas/nudge_sms"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/message"
)

func TestNudgeSMSLambdaHandler_HandleRequest(t *testing.T) {
	// test.SetEnvVars()

	// Setup mock database
	// db, _, err := test.SetupMockDB()
	// require.NoError(t, err, "Could not set up mock db")

	database := db.Get(config.Get())

	memSvc := message.NewMemoryService(
		message.NewMessageRepository(database), 1,
	)

	llmSvc := &ai.OpenAICompletionService{}

	handler := &main.NudgeSMSLambdaHandler{
		MemoryService:          memSvc,
		CompletionService:      llmSvc,
		MaxMemories:            main.MaxMemories,
		NudgeIfNoMessagesSince: main.TimeSinceLastMessage,
	}

	handler.Init(database)

	testUser := &models.User{
		ID:          36,
		PhoneNumber: "+12533243071",
		Firstname:   "Kevin",
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	err := handler.Nudge(testUser, wg)
	wg.Wait()

	require.NoError(t, err, "There should be no errors when nudging.")
}
