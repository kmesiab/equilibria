package transaction_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/test"
	"github.com/kmesiab/equilibria/lambdas/lib/transaction"
	"github.com/kmesiab/equilibria/lambdas/models"
)

var (
	TestTransactionColumns = sqlmock.NewRows([]string{
		"id",
		"user_id",
		"conversation_id",
		"transaction_type",
		"funding_source",
		"description",
		"reference_id",
		"timestamp",
	})

	now = time.Now()

	TestTransaction = models.Transaction{
		ID:              1,
		UserID:          1,
		ConversationID:  1,
		Amount:          10,
		TransactionType: "credit",
		FundingSource:   "stripe",
		Description:     "Unit Test - FindByID",
		ReferenceID:     "xxx",
		Timestamp:       &now,
	}
)

func TestTransactionRepository_FindByID(t *testing.T) {
	test.SetEnvVars()
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := transaction.NewTransactionRepository(db)

	mock.ExpectQuery("SELECT \\* FROM `transactions`").
		WithArgs(1).WillReturnRows(
		TestTransactionColumns.AddRow(
			1, 1, 1,
			TestTransaction.TransactionType,
			TestTransaction.FundingSource,
			TestTransaction.Description,
			TestTransaction.ReferenceID,
			time.Now(),
		))

	mock.ExpectQuery("SELECT \\* FROM `conversations`").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).
			AddRow(1, 1))

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
		AddRow(1, "jane@example.com"))

	mock.ExpectQuery("SELECT \\* FROM `users`").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
		AddRow(1, "jane@example.com"))

	txn, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, txn)
	assert.Equal(t, int64(1), txn.ID)
}

func TestTransactionRepository_Create(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := transaction.NewTransactionRepository(db)

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO `transactions`").
		WithArgs(
			sqlmock.AnyArg(),       // for created_at
			sqlmock.AnyArg(),       // for updated_at
			sqlmock.AnyArg(),       // for deleted_at
			1,                      // user_id
			1,                      // conversation_id
			float64(10),            // amount
			"credit",               // transaction_type
			"stripe",               // funding_source
			"Unit Test - FindByID", // description
			"xxx",                  // reference_id
			1,                      // id (if it's not auto-incremented or if you're explicitly setting it)
			sqlmock.AnyArg(),       // timestamp
		).WillReturnResult(test.GenerateMockLastAffectedRow())

	mock.ExpectCommit()

	err = repo.Create(&TestTransaction)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), TestTransaction.ID)
}

func TestTransactionRepository_Update(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	now := time.Now()
	repo := transaction.NewTransactionRepository(db)

	TestTransaction := &models.Transaction{
		ID:              1,
		UserID:          1,
		ConversationID:  1,
		Amount:          10,
		TransactionType: "credit",
		FundingSource:   "stripe",
		Description:     "Unit Test - FindByID",
		ReferenceID:     "xxx",
		Timestamp:       &now,
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `transactions` SET").
		WithArgs(
			sqlmock.AnyArg(),                // for updated_at
			TestTransaction.UserID,          // user_id
			TestTransaction.ConversationID,  // conversation_id
			float64(10),                     // amount
			TestTransaction.TransactionType, // transaction_type
			TestTransaction.FundingSource,   // funding_source
			TestTransaction.Description,     // description
			TestTransaction.ReferenceID,     // reference_id
			sqlmock.AnyArg(),                // timestamp
			1,                               // WHERE id = ?
		).WillReturnResult(sqlmock.NewResult(0, 1)) // Assuming one row is affected
	mock.ExpectCommit()

	err = repo.Update(TestTransaction)

	assert.NoError(t, err)
	assert.NotNil(t, TestTransaction)
	assert.Equal(t, int64(1), TestTransaction.ID)
}

func TestTransactionRepository_Debit(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := transaction.NewTransactionRepository(db)

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO `transactions`").
		WithArgs(
			sqlmock.AnyArg(),            // for created_at
			sqlmock.AnyArg(),            // for updated_at
			sqlmock.AnyArg(),            // for deleted_at
			1,                           // user_id
			1,                           // conversation_id
			-10.500001,                  // amount, enforced negatives
			"debit",                     // transaction_type
			"stripe",                    // funding_source
			TestTransaction.Description, // description
			"xxx",                       // reference_id
		).WillReturnResult(test.GenerateMockLastAffectedRow())

	mock.ExpectCommit()

	txn, err := repo.Debit(
		int64(1),
		int64(1),
		10.500001, // maintain the precision
		models.FundingSourceStringStripe,
		TestTransaction.ReferenceID,
		TestTransaction.Description,
	)

	assert.NoError(t, err)
	assert.NotNil(t, txn)
	assert.Negativef(t, txn.Amount, "Amount should be converted to a negative")
	assert.Equal(t, int64(1), TestTransaction.ID)
}

func TestTransactionRepository_Credit(t *testing.T) {
	db, mock, err := test.SetupMockDB()
	assert.NoError(t, err)

	repo := transaction.NewTransactionRepository(db)

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO `transactions`").
		WithArgs(
			sqlmock.AnyArg(),                   // for created_at
			sqlmock.AnyArg(),                   // for updated_at
			sqlmock.AnyArg(),                   // for deleted_at
			1,                                  // user_id
			1,                                  // conversation_id
			10.500001,                          // amount, enforced negatives
			models.TransactionTypeStringCredit, // transaction_type
			models.FundingSourceStringStripe,   // funding_source
			TestTransaction.Description,        // description
			TestTransaction.ReferenceID,        // reference_id
		).WillReturnResult(test.GenerateMockLastAffectedRow())

	mock.ExpectCommit()

	txn, err := repo.Credit(
		int64(1),
		int64(1),
		-10.500001, // maintain the precision
		models.FundingSourceStringStripe,
		TestTransaction.ReferenceID,
		TestTransaction.Description,
	)

	assert.NoError(t, err)
	assert.NotNil(t, txn)
	assert.Positive(t, txn.Amount, "Amount should be converted to a negative")
	assert.Equal(t, int64(1), TestTransaction.ID)
}
