package message_type_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/message_type"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/models"
)

func TestMain(m *testing.M) {
	test.SetEnvVars()

	code := m.Run()

	os.Exit(code)
}

func TestMessageTypeRepository_FindByID(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := message_type.NewMessageTypeRepository(db)

	var msgType = models.NewMessageTypeSMS()

	mock.ExpectQuery("SELECT \\* FROM `message_types`").
		WithArgs(msgType.ID).
		WillReturnRows(test.GenerateMockMessageType())

	messageType, err := repo.FindByID(2)

	assert.NoError(t, err)
	assert.NotNil(t, messageType)

	assert.Equal(t, msgType.ID, messageType.ID)
	assert.Equal(t, msgType.Name, messageType.Name)
	assert.Equal(t, msgType.BillRateInCredits, messageType.BillRateInCredits)
}

func TestMessageTypeRepository_FindByName(t *testing.T) {

	var (
		expectedMessageTypeID     = int64(2)
		expectedBillRateInCredits = 0.5
	)

	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT \\* FROM `message_types`").
		WithArgs("Email").WillReturnRows(test.GenerateMockMessageType())

	repo := message_type.NewMessageTypeRepository(db)
	messageType, err := repo.FindByName("Email")

	assert.NoError(t, err, "There shouldn't be any errors when finding by name")
	assert.NotNil(t, messageType, "The message type shouldn't be nil")
	assert.Equal(t, expectedMessageTypeID, messageType.ID, "The message type should have an ID of 1")
	assert.Equal(t, "SMS", messageType.Name, "The message type should be Email")
	assert.Equal(t, expectedBillRateInCredits, messageType.BillRateInCredits, "The message type should have a bill rate of 0")
}

func TestMessageTypeRepository_GetAll(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	rows := test.GenerateMockMessageType()
	rows.AddRow(1, "Email", 0.5)

	mock.ExpectQuery("SELECT \\* FROM `message_types`").
		WillReturnRows(rows)

	repo := message_type.NewMessageTypeRepository(db)
	messageTypes, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, messageTypes, 2)
	assert.Equal(t, "SMS", messageTypes[0].Name)
	assert.Equal(t, 0.5, messageTypes[0].BillRateInCredits)
	assert.Equal(t, "Email", messageTypes[1].Name)
	assert.Equal(t, 0.5, messageTypes[1].BillRateInCredits)
}
