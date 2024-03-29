package message_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/lib/message"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/models"
)

func TestMessageRepository_FindByUser(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := message.NewMessageRepository(db)

	mock.ExpectQuery(test.MessageSelectQuery).WithArgs(1).
		WillReturnRows(test.GenerateMockMessageRepositoryMessages())

	mock.ExpectQuery("SELECT \\* FROM `conversations`").WithArgs(1).
		WillReturnRows(test.GenerateMockConversation(false))

	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectMessageStatusAndTypes(&mock)

	test.ExpectMockSelectUser(&mock, 1)

	messages, err := repo.FindByUser(&models.User{ID: 1})

	assert.NoError(t, err,
		"Error should be nil when getting message by ID")

	assert.NotNil(t, messages, "Message should not be nil")

	assert.Greater(t, len(*messages), 0, "Should have found at least one message")

	// take first user
	msg := (*messages)[0]

	assert.Greater(t, msg.ID, int64(0),
		"User ID should be populated")

	// Check for preloaded relationships
	assert.Greater(t, msg.ID, int64(0),
		"The conversation's user should be set")

	assert.Greater(t, msg.Conversation.User.ID, int64(0),
		"The conversation's user should be set")

	assert.Greater(t, msg.From.ID, int64(0),
		"The from user should be set")

	assert.Greater(t, msg.To.ID, int64(0),
		"The from to should be set")

	assert.Greater(t, msg.MessageType.ID, int64(0),
		"The message type ID to should be set")

	assert.Greater(t, msg.MessageStatus.ID, int64(0),
		"The message status ID to should be set")

}

func TestMessageRepository_FindByID(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := message.NewMessageRepository(db)

	mock.ExpectQuery(test.MessageSelectQuery).WithArgs(1, sqlmock.AnyArg()).
		WillReturnRows(test.GenerateMockMessageRepositoryMessages())

	mock.ExpectQuery("SELECT \\* FROM `conversations`").WithArgs(1).
		WithArgs(1).WillReturnRows(test.GenerateMockConversation(false))

	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectUser(&mock, 1)
	test.ExpectMockSelectMessageStatusAndTypes(&mock)
	test.ExpectMockSelectUser(&mock, 1)

	msg, err := repo.FindByID(1)

	assert.NoError(t, err,
		"Error should be nil when getting message by ID")

	assert.NotNil(t, msg, "Message should not be nil")

	assert.Greater(t, msg.ID, int64(0),
		"User ID should be populated")

	// Check for preloaded relationships
	assert.Greater(t, msg.Conversation.ID, int64(0),
		"The conversation's user should be set")

	assert.Greater(t, msg.Conversation.User.ID, int64(0),
		"The conversation's user should be set")

	assert.Greater(t, msg.From.ID, int64(0),
		"The from user should be set")

	assert.Greater(t, msg.To.ID, int64(0),
		"The from to should be set")

	assert.Greater(t, msg.MessageType.ID, int64(0),
		"The message type ID to should be set")

	assert.Greater(t, msg.MessageStatus.ID, int64(0),
		"The message status ID to should be set")

}

func TestMessageRepository_GetNoRows(t *testing.T) {

	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := message.NewMessageRepository(db)

	mock.ExpectQuery(test.MessageSelectQuery).WithArgs(1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err = repo.FindByID(1)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestGetMessageDatabaseError(t *testing.T) {

	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := message.NewMessageRepository(db)

	mock.ExpectQuery(test.MessageSelectQuery).WillReturnError(errors.New("db error"))

	_, err = repo.FindByID(1)

	assert.Error(t, err)
}

func TestMessageRepository_Create(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := message.NewMessageRepository(db)

	now := time.Now()

	msg := &models.Message{
		ConversationID:  int64(1),
		FromUserID:      int64(1),
		ToUserID:        int64(2),
		SentAt:          &now,
		Body:            "Test *outbound* message content",
		MessageTypeID:   1,
		MessageStatusID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `messages`").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Create(msg)
	assert.NoError(t, err)

	assert.Greater(t, msg.ID, int64(0),
		"Message ID should be populated")

}

func TestMessageRepository_Update(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := message.NewMessageRepository(db)

	msg := &models.Message{
		ID:   1,
		Body: "Updated message content",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `messages` SET `body`=").
		WithArgs("Updated message content", sqlmock.AnyArg(), 1).
		WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Update(msg)
	assert.NoError(t, err)
	assert.NotEqual(t, msg.UpdatedAt, time.Time{}, "Updated at should be populated")
}

func TestMessageRepository_DeleteIsSoftDelete(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	require.NoError(t, err)

	repo := message.NewMessageRepository(db)

	msg := &models.Message{
		ID:   1,
		Body: "Updated message content",
	}

	mock.ExpectBegin()
	mock.ExpectExec(
		"UPDATE `messages` SET `body`=\\?,`updated_at`=\\? WHERE `messages`.`deleted_at` IS NULL AND `id` = \\?",
	).WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), 1,
	).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Update(msg)
	assert.NoError(t, err)
}
