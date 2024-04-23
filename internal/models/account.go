package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id" validate:"required" gorm:"primaryKey"`
	Balance   float64   `json:"balance" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
