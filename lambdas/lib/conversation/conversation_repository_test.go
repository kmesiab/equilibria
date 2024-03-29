package conversation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/conversation"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/models"
)

var now = time.Now()

func TestConversationRepository_FindByID(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(1, sqlmock.AnyArg()).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	convo, err := repo.FindByID(1)

	assert.NoError(t, err,
		"Error should be nil when getting conversation by ID")

	assert.NotNil(t, convo, "Conversation should not be nil")

	assert.Greater(t, convo.User.ID, int64(0),
		"User ID should be populated")

	assert.Greater(t, convo.User.AccountStatus.ID, int64(0),
		"Account status should be populated")
}

func TestConversationRepository_GetAll(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectQuery(test.ConversationSelectQuery).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(2, sqlmock.AnyArg()).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	conversations, err := repo.GetAll()

	assert.NoError(t, err,
		"Error should be nil when getting conversation by ID")

	assert.NotNil(t, conversations, "Conversation should not be nil")

	assert.Greater(t, len(conversations), 0, "Conversation slice should not be empty")

	convo := conversations[0]

	assert.Greater(t, convo.User.ID, int64(0),
		"User ID should be populated")

	assert.Greater(t, convo.User.AccountStatus.ID, int64(0),
		"Account status should be populated")
}

func TestConversationRepository_FindByUser(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(2).
		WillReturnRows(test.GenerateMockConversation(false))

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockAccountStatusPending())

	convos, err := repo.FindByUser(models.User{ID: 1})

	require.NoError(t, err)
	assert.Greater(t, len(*convos), 0)

	convo := (*convos)[0]

	assert.NotEqual(t, convo.User.AccountStatus.ID, 0,
		"User account status should be set")

	assert.NotEqual(t, convo.User.ID, 0,
		"User should be set")

}

func TestConversationRepository_NoResults(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectQuery(test.ConversationSelectQuery).WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{}))

	convos, err := repo.FindByUser(models.User{ID: 99})

	assert.NoError(t, err)
	assert.Equal(t, 0, len(*convos))
}

func TestConversationRepository_DatabaseError(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectQuery(test.ConversationSelectQuery).WillReturnError(errors.New("db error"))

	_, err = repo.FindByUser(models.User{ID: 1})

	assert.Error(t, err)
}

func TestConversationRepository_Create(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	convo := &models.Conversation{
		UserID:    2,
		StartTime: &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `conversations`").WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Create(convo)
	assert.NoError(t, err)

	assert.Greater(t, convo.ID, int64(0),
		"Conversation ID should be populated")

	assert.Nil(t, convo.EndTime, "End time should be nil")

}

func TestConversationRepository_Update(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	convo := &models.Conversation{
		ID:      1,
		UserID:  1,
		EndTime: &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `conversations`").WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Update(convo)
	assert.NoError(t, err)
}

func TestConversationRepository_SoftDelete(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `conversations`").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.SoftDelete(1)
	assert.NoError(t, err)
}

func TestConversationRepository_HardDelete(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `conversations` WHERE id = ?").
		WithArgs(1).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.HardDelete(1)
	assert.NoError(t, err)
}

func TestConversationRepository_StartConversation(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `conversations` SET `start_time`=\\?,`updated_at`=\\? WHERE `id` = \\?").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.StartConversation(1)
	assert.NoError(t, err)
}

func TestConversationRepository_EndConversation(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `conversations`").WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.EndConversation(1)
	assert.NoError(t, err)
}

func TestConversationRepository_GetOpenConvByUserIDMultipleConvos(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	// Mock query for fetching open conversations
	mock.ExpectQuery("SELECT \\* FROM `conversations`").
		WithArgs(1).
		WillReturnRows(test.GenerateMockOpenConversations())

	// Mock query for fetching user and account status
	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())
	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	conversations, err := repo.GetOpenConversationsByUserID(1)

	require.NoError(t, err)
	assert.NotNil(t, conversations)
	assert.Greater(t, len(*conversations), 1, "There should be at least one open conversation")

	convo := (*conversations)[0]
	assert.Equal(t, int64(1), convo.UserID, "User ID should match")
	assert.Nil(t, convo.EndTime, "EndTime should be nil for open conversations")
}

func TestConversationRepository_GetOpenConversationsByUserID(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	// Mock query for fetching open conversations
	mock.ExpectQuery("SELECT \\* FROM `conversations`").
		WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	// Mock query for fetching user and account status
	test.ExpectMockSelectUser(&mock, 1)

	conversations, err := repo.GetOpenConversationsByUserID(1)

	require.NoError(t, err)
	assert.NotNil(t, conversations)
	assert.Greater(t, len(*conversations), 0, "There should be at least one open conversation")

	convo := (*conversations)[0]
	assert.Equal(t, int64(1), convo.UserID, "User ID should match")
	assert.Nil(t, convo.EndTime, "EndTime should be nil for open conversations")
}

func TestConversationRepository_GetOpenConversationsByUserIDNoResults(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := conversation.NewConversationRepository(db)

	// Mock query for fetching open conversations
	mock.ExpectQuery("SELECT \\* FROM `conversations`").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(nil))

	// Mock query for fetching user and account status
	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1, sqlmock.AnyArg()).
		WillReturnRows(test.GenerateMockUserRepositoryUser())
	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").WithArgs(1).
		WillReturnRows(test.GenerateMockUserRepositoryUser())

	conversations, err := repo.GetOpenConversationsByUserID(1)

	require.NoError(t, err)
	assert.NotNil(t, conversations)
	assert.Equal(t, len(*conversations), 0, "There should be no conversations")
}
