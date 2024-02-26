package user_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/hasher"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/lib/user"
	"github.com/kmesiab/equilibria/lambdas/models"
)

func TestUserRepository_Create(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)

	pwd := test.DefaultTestPassword

	newUser := models.User{
		PhoneNumber:     test.DefaultTestPhoneNumber,
		PhoneVerified:   true,
		Firstname:       test.DefaultTestUserFirstname,
		Lastname:        test.DefaultTestUserLastname,
		Email:           test.DefaultTestEmail,
		Password:        &pwd,
		AccountStatusID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(
			test.DefaultTestPassword,
			newUser.PhoneNumber,
			newUser.PhoneVerified,
			newUser.Firstname,
			newUser.Lastname,
			newUser.Email,
			sqlmock.AnyArg(), // For CreatedAt
		).WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	mock.ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs("Pending Activation").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Pending Activation"))
	err = repo.Create(&newUser)

	assert.Equal(t, int64(1), newUser.AccountStatus.ID)
	assert.Equal(t, int64(1), newUser.AccountStatusID)
	assert.NoError(t, err)
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)

	newUser := models.User{
		ID:    1,
		Email: "john@example.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(test.GenerateMockLastAffectedRow())

	mock.ExpectCommit()

	err = repo.Update(&newUser)
	assert.NoError(t, err)
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)
	userId := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `users`").
		WithArgs(userId).
		WillReturnResult(test.GenerateMockLastAffectedRow())
	mock.ExpectCommit()

	err = repo.Delete(userId)
	assert.NoError(t, err)
}

func TestUserRepository_ListAll(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(1, "jane@example.com").
			AddRow(2, "john@example.com"))

	users, err := repo.ListAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepository_FindByEmail(t *testing.T) {

	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)
	email := "jane@example.com"
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = ?").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(1, email))

	newUser, err := repo.FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, email, newUser.Email)
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)
	id := int64(1)
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id`").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(1, id))

	newUser, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, id, newUser.ID)
}

func TestUserRepository_FindByPhoneNumber(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)
	phoneNumber := "jane@example.com"
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE phone_number =").
		WithArgs(phoneNumber).
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number"}).
			AddRow(1, phoneNumber))

	newUser, err := repo.FindByPhoneNumber(phoneNumber)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, phoneNumber, newUser.PhoneNumber)
}

func TestUserRepository_FindByName(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)
	firstName := "jane"

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE firstname LIKE \\? OR lastname LIKE \\? ").
		WithArgs("%"+firstName+"%", "%"+firstName+"%"). // Two separate arguments for firstname and lastname
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname"}).
			AddRow(1, firstName))

	newUser, err := repo.FindByName(firstName)
	assert.NoError(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, firstName, newUser.Firstname)
}

func TestUserRepository_CheckPassword(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := user.NewUserRepository(db)

	password := "testPassword"
	hashedPassword, err := hasher.HashPassword(password)

	require.NoError(t, err,
		"There should be no errors when hashing the password.")

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE phone_number =").
		WithArgs(test.DefaultTestPhoneNumber).
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "password"}).
			AddRow(test.DefaultFromUserID, test.DefaultTestPhoneNumber, hashedPassword))

	valid, err := repo.CheckPassword(test.DefaultTestPhoneNumber, password)
	assert.NoError(t, err)
	assert.True(t, valid)
}
