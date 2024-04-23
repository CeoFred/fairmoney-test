package acme_http_mock

import "github.com/gofrs/uuid"

type Transaction struct {
	Reference uuid.UUID `json:"reference" validate:"required"`
	Amount    float64   `json:"amount" validate:"required"`
	AccountID uuid.UUID `json:"account_id" validate:"required"`
}
