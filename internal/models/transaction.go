package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type TransactionStatus string
type TransactionIntent string

const (
	PendingTransaction    TransactionStatus = "Pending"
	SuccessfulTransaction TransactionStatus = "Successful"
	FailedTransaction     TransactionStatus = "Failed"
)

const (
	Credit TransactionIntent = "Credit"
	Debit  TransactionIntent = "Debit"
)

type Transaction struct {
	ID        uuid.UUID         `json:"id" gorm:"primaryKey" validate:"required"`
	Amount    float64           `json:"amount" validate:"required"`
	Type      TransactionIntent `json:"type" validate:"required"`
	AccountID uuid.UUID         `json:"account_id" validate:"required"`
	Status    TransactionStatus `json:"status" validate:"required"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
