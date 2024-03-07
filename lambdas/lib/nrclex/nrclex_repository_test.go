package nrclex

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestNrcLexRepository_Create(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	now := time.Now()

	nrcLex := &models.NrcLex{
		UserID:    1,
		MessageID: 2,
		Anger:     0.5,
		Joy:       0.5,
		CreatedAt: now,
	}

	mock.ExpectBegin() // Expect transaction begin
	mock.ExpectExec("INSERT INTO `nrclex`").WithArgs(
		1,                // user_id
		2,                // message_id
		0.5,              // anger
		0.0,              // anticipation
		0.0,              // disgust
		0.0,              // fear
		0.0,              // trust
		0.5,              // joy
		0.0,              // negative
		0.0,              // positive
		0.0,              // sadness
		0.0,              // surprise
		0.0,              // vader_compound
		0.0,              // vader_neg
		0.0,              // vader_neu
		0.0,              // vader_pos
		sqlmock.AnyArg(), // created_at
		sqlmock.AnyArg(), // updated_at
		sqlmock.AnyArg(), // deleted_at (can be nil or a specific time)
	).WillReturnResult(sqlmock.NewResult(1, 1)) // Assuming ID 1 is returned after insert
	mock.ExpectCommit() // Expect transaction commit

	err = repo.Create(nrcLex)
	assert.NoError(t, err)
}

func TestNrcLexRepository_FindByID(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `nrclex` WHERE").
		WithArgs(
			1, // the actual id of the NrcLex
			1, // this is because .First() set's limit=1
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "message_id", "anger", "joy", "created_at",
		}).AddRow(1, 1, 2, 0.5, 0.5, time.Now()))

	nrcLex, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, nrcLex)
}

func TestNrcLexRepository_Update(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	nrcLex := &models.NrcLex{
		UserID:    1,
		MessageID: 2,
		ID:        1,
		Anger:     0.6, // Updated value
		Joy:       0.4, // Updated value
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `nrclex`").
		WithArgs(
			1,
			2,
			0.6,              // anger
			0.0,              // anticipation
			0.0,              // disgust
			0.0,              // fear
			0.0,              // trust
			0.4,              // joy
			0.0,              // negative
			0.0,              // positive
			0.0,              // sadness
			0.0,              // surprise
			0.0,              // vader_compound
			0.0,              // vader_neg
			0.0,              // vader_neu
			0.0,              // vader_pos
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at, expecting a dynamic value (the current timestamp in this case).
			nil,              // deleted_at is NULL.
			1,                // WHERE condition, matching the record to update by ID.
		).WillReturnResult(sqlmock.NewResult(0, 1)) // Assuming 1 row affected.
	mock.ExpectCommit()

	err = repo.Update(nrcLex)
	assert.NoError(t, err)
}

func TestNrcLexRepository_Delete(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `nrclex` WHERE").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(1)
	assert.NoError(t, err)
}

func TestFindNrcLex_NoRows(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `nrclex` WHERE").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	_, err = repo.FindByID(1)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestNrcLexRepository_Create_FKViolation(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	// Attempt to create an NrcLex with IDs that could violate FK constraints
	nrcLex := &models.NrcLex{
		UserID:    999, // Assuming this ID does not exist
		MessageID: 999, // Assuming this ID does not exist
		Anger:     0.5,
		Joy:       0.5,
		CreatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `nrclex`").
		WithArgs(
			nrcLex.UserID,
			nrcLex.MessageID,
			nrcLex.Anger,
			0.0, 0.0, 0.0, 0.0, // Other emotions set to default
			nrcLex.Joy,
			0.0, 0.0, 0.0, 0.0, // Other sentiments set to default
			0.0, 0.0, 0.0, 0.0, // Vader scores set to default
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("foreign key constraint fails"))
	mock.ExpectRollback()

	err = repo.Create(nrcLex)
	assert.Error(t, err)
}

func TestNrcLexRepository_Update_NonExistentID(t *testing.T) {
	db, mock, err := setupMockDB()
	require.NoError(t, err)

	repo := NewRepository(db)

	// Assuming nrcLex is the struct for a record that does not exist in the database.
	nrcLex := &models.NrcLex{
		ID:        99999,      // Use an ID that does not exist.
		UserID:    1,          // Dummy data for the sake of completeness.
		MessageID: 2,          // Dummy data for the sake of completeness.
		Anger:     0.6,        // Updated value
		Joy:       0.4,        // Updated value
		CreatedAt: time.Now(), // Dummy timestamp
		UpdatedAt: time.Now(), // Dummy timestamp
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `nrclex`").
		WithArgs(
			1,
			2,
			0.6,              // anger
			0.0,              // anticipation
			0.0,              // disgust
			0.0,              // fear
			0.0,              // trust
			0.4,              // joy
			0.0,              // negative
			0.0,              // positive
			0.0,              // sadness
			0.0,              // surprise
			0.0,              // vader_compound
			0.0,              // vader_neg
			0.0,              // vader_neu
			0.0,              // vader_pos
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at, expecting a dynamic value (the current timestamp in this case).
			nil,              // deleted_at is NULL.
			99999,            // WHERE condition, matching the record to update by ID.
		).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `nrclex`")

	err = repo.Update(nrcLex)
	assert.Error(t, err)
}
