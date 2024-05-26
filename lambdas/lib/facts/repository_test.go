package facts_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/facts"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/models"
)

func TestFactsRepository_Create(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)

	newFact := models.Fact{
		UserID:         1,
		ConversationID: 1,
		Body:           "Test fact",
		Reasoning:      "Test reasoning",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `facts`").
		WithArgs(
			newFact.UserID,
			newFact.ConversationID,
			newFact.Body,
			newFact.Reasoning,
		).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Create(&newFact)
	assert.NoError(t, err)
}

func TestFactsRepository_Update(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)

	newFact := models.Fact{
		ID:             1,
		UserID:         1,
		ConversationID: 1,
		Body:           "Updated fact",
		Reasoning:      "Updated reasoning",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `facts` SET").
		WithArgs(
			newFact.UserID,
			newFact.ConversationID,
			newFact.Body,
			newFact.Reasoning,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			newFact.ID,
		).WillReturnResult(test.GenerateMockLastAffectedRow())

	mock.ExpectCommit()

	err = repo.Update(&newFact)
	assert.NoError(t, err)
}

func TestFactsRepository_Delete(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)
	factId := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `facts` WHERE `facts`.`id` = ?").
		WithArgs(factId).
		WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Delete(factId)
	assert.NoError(t, err)
}

func TestFactsRepository_FindByID(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)
	id := int64(1)
	mock.ExpectQuery("SELECT \\* FROM `facts` WHERE `facts`.`id` = ?").
		WithArgs(id, 1). // Extra "1" for the limit = 1
		WillReturnRows(sqlmock.NewRows([]string{"id", "body"}).
			AddRow(1, "Fact 1"))

	newFact, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, newFact)
	assert.Equal(t, id, newFact.ID)
}

func TestFactsRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)
	id := int64(1)
	mock.ExpectQuery("SELECT \\* FROM `facts` WHERE `facts`.`id` = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "body"}))

	newFact, err := repo.FindByID(id)
	assert.Error(t, err)
	assert.Nil(t, newFact)
}

func TestFactsRepository_Create_Error(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)

	newFact := models.Fact{
		UserID:         1,
		ConversationID: 1,
		Body:           "Test fact",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `facts`").
		WithArgs(
			newFact.UserID,
			newFact.ConversationID,
			newFact.Body,
		).WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.Create(&newFact)
	assert.Error(t, err)
}

func TestFactsRepository_Update_Error(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)

	newFact := models.Fact{
		ID:             1,
		UserID:         1,
		ConversationID: 1,
		Body:           "Updated fact",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `facts` SET").
		WithArgs(
			newFact.UserID,
			newFact.ConversationID,
			newFact.Body,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			newFact.ID,
		).WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.Update(&newFact)
	assert.Error(t, err)
}

func TestFactsRepository_Delete_Error(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := facts.NewRepository(db)
	factId := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `facts` WHERE `facts`.`id` = ?").
		WithArgs(factId).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err = repo.Delete(factId)
	assert.Error(t, err)
}
