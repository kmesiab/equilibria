package status_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/status"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

var (
	TestStatusColumns = sqlmock.NewRows([]string{"id", "name"})
)

func TestStatusRepository_FindByID(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := status.NewStatusRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs(1, sqlmock.AnyArg()).WillReturnRows(
		TestStatusColumns.AddRow(1, "Active"))

	accountStatus, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, accountStatus)
	assert.Equal(t, int64(1), accountStatus.ID)
	assert.Equal(t, "Active", accountStatus.Name)
}

func TestStatusRepository_FindByName(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := status.NewStatusRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs("Active", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Active"))

	accountStatus, err := repo.FindByName("Active")

	assert.NoError(t, err)
	assert.NotNil(t, accountStatus)
	assert.Equal(t, int64(1), accountStatus.ID)
	assert.Equal(t, "Active", accountStatus.Name)
}

func TestStatusRepository_GetAll(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := status.NewStatusRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Active").
			AddRow(2, "Suspended"))

	accountStatuses, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, accountStatuses, 2)
	assert.Equal(t, "Active", accountStatuses[0].Name)
	assert.Equal(t, "Suspended", accountStatuses[1].Name)
}
