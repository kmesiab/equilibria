package transaction

import (
	"math"

	"gorm.io/gorm"

	"github.com/kmesiab/equilibria/lambdas/models"
)

// TransactionRepository is a repository for managing Transactions.
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new TransactionRepository.
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {

	return &TransactionRepository{
		db: db,
	}
}

// Create creates a new transaction in the database.
func (repo *TransactionRepository) Create(transaction *models.Transaction) error {

	err := repo.db.Create(transaction).Error

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

// Update updates a transaction's details in the database.
func (repo *TransactionRepository) Update(transaction *models.Transaction) error {

	return repo.db.Updates(transaction).Error
}

// Delete deletes a transaction from the database.
func (repo *TransactionRepository) Delete(id int64) error {
	return repo.db.Delete(&models.Transaction{}, id).Error
}

// FindByID finds a transaction by their ID.
func (repo *TransactionRepository) FindByID(id int64) (*models.Transaction, error) {
	var transaction models.Transaction
	err := repo.db.Preload("User.AccountStatus").
		Preload("Conversation.User.AccountStatus").
		First(&transaction, id).
		Error

	return &transaction, err
}

// Credit creates a new credit transaction in the database.
func (repo *TransactionRepository) Credit(
	userID int64,
	conversationID int64,
	amount float64,
	fundingSource, referenceID, Description string,
) (*models.Transaction, error) {
	transaction := models.Transaction{
		UserID:          userID,
		ConversationID:  conversationID,
		Amount:          amount,
		TransactionType: models.TransactionTypeStringCredit,
		FundingSource:   fundingSource,
		Description:     Description,
		ReferenceID:     referenceID,
	}

	transaction.Amount = math.Abs(transaction.Amount)

	err := repo.Create(&transaction)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// Debit creates a new credit transaction in the database.
func (repo *TransactionRepository) Debit(
	userID int64,
	conversationID int64,
	amount float64,
	fundingSource, referenceID, Description string,
) (*models.Transaction, error) {
	transaction := models.Transaction{
		UserID:          userID,
		ConversationID:  conversationID,
		Amount:          amount,
		TransactionType: models.TransactionTypeStringDebit,
		FundingSource:   fundingSource,
		Description:     Description,
		ReferenceID:     referenceID,
	}

	transaction.Amount = math.Abs(transaction.Amount) * -1

	err := repo.db.Create(&transaction).Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
