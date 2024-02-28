package message_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/models"
)

func TestMemoryService_GeneratePrompt(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery(test.MessageSelectQuery).WithArgs(1).
		WillReturnRows(test.GenerateMockMessageRepositoryMessages())

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	mock.ExpectQuery("SELECT \\* FROM `message_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockMessageStatus())

	mock.ExpectQuery("SELECT \\* FROM `message_types`").WithArgs(
		models.NewMessageTypeSMS().ID,
	).WillReturnRows(test.GenerateMockMessageType())

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	mock.ExpectQuery("SELECT \\* FROM `message_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockMessageStatus())

	repo := message.NewMemoryService(message.NewMessageRepository(db), 1)

	prompt, err := repo.GeneratePrompt(&models.User{ID: 1})

	require.NoError(t, err)
	require.NotEmpty(t, prompt)

}
