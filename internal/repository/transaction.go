package repository

import (
	"errors"
	"fmt"

	"github.com/CeoFred/fairmoney/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	database *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		database: db,
	}
}

func (a TransactionRepository) FindRecordsByCondition(condition, value string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := a.database.Raw(fmt.Sprintf(`SELECT * FROM transactions WHERE %s = ?`, condition), value).Scan(&transactions).Error
	if err != nil {
		return nil, err
	}
	if transactions != nil {
		return transactions, nil
	}
	return nil, nil
}

func (a TransactionRepository) Find(id uuid.UUID) (*models.Transaction, error) {
	var transaction *models.Transaction
	err := a.database.First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (a TransactionRepository) Exists(id uuid.UUID) (bool, error) {
	var transaction *models.Transaction
	err := a.database.First(&transaction, "id = ?", id).Error
	if err != nil {
		return false, err
	}

	if transaction == nil {
		return false, nil
	}
	return true, nil
}

func (a TransactionRepository) Create(transaction *models.Transaction) error {
	return a.database.Model(&models.Transaction{}).Create(transaction).Error
}

func (a TransactionRepository) Save(transaction *models.Transaction) (*models.Transaction, error) {

	txn := a.database.Model(transaction).Where("id = ?", transaction.ID).Updates(&transaction).First(transaction)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return transaction, nil
}
