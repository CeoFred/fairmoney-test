package repository

import (
	"errors"
	"fmt"

	"github.com/CeoFred/fairmoney/internal/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	database *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		database: db,
	}
}

func (a *AccountRepository) FindRecordsByCondition(condition, value string) ([]*models.Account, error) {
	var accounts []*models.Account
	err := a.database.Raw(fmt.Sprintf(`SELECT * FROM accounts WHERE %s = ?`, condition), value).Scan(&accounts).Error
	if err != nil {
		return nil, err
	}
	if accounts != nil {
		return accounts, nil
	}
	return nil, nil
}

func (a *AccountRepository) Find(id uuid.UUID) (*models.Account, error) {
	var account *models.Account
	err := a.database.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *AccountRepository) Create(account *models.Account) error {
	return a.database.Model(&models.Account{}).Create(account).Error
}

func (a *AccountRepository) Save(account *models.Account) (*models.Account, error) {

	txn := a.database.Save(account)

	if txn.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}

	if txn.Error != nil {
		return nil, txn.Error
	}

	return account, nil
}
