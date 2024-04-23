package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/CeoFred/fairmoney/database"
	"github.com/CeoFred/fairmoney/internal/handlers"
	"github.com/CeoFred/fairmoney/internal/models"
	"github.com/CeoFred/fairmoney/internal/repository"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
type SuccessResponse struct {
	Data    models.Transaction `json:"data"`
	Message string             `json:"message"`
	Status  bool               `json:"status"`
}

var accountID = uuid.FromStringOrNil("018ec333-5c51-7f5d-b3fc-218d742e9a02")

func TestRouter(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "{\"message\":\"Whooops! Not Found\"}", w.Body.String())
}

func TestInvalidAccountOnTransaction(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	transactionJSON, err := json.Marshal(handlers.TransactinRequest{
		AccountID: uuid.Must(uuid.NewV4()),
		Reference: uuid.Must(uuid.NewV4()),
		Amount:    400,
		Type:      models.Debit,
	})

	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(transactionJSON))
	router.ServeHTTP(w, req)
	var response ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	assert.Equal(t, "record not found", response.Error)
}

func TestValidTransaction(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()

	transactionJSON, err := json.Marshal(handlers.TransactinRequest{
		AccountID: accountID,
		Reference: uuid.Must(uuid.NewV4()),
		Amount:    400,
		Type:      models.Credit,
	})

	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(transactionJSON))
	router.ServeHTTP(w, req)
	var response SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, response.Status, true)
}

func TestNegativeBalanceTransaction(t *testing.T) {
	router := SetupRouter()
	sw := httptest.NewRecorder()

	accountRepo := repository.NewAccountRepository(database.DB)
	_, err := accountRepo.Save(&models.Account{
		ID:      accountID,
		Balance: 0,
	})
	if err != nil {
		t.Error(err)
	}

	transactionJSON, err := json.Marshal(handlers.TransactinRequest{
		AccountID: accountID,
		Reference: uuid.Must(uuid.NewV4()),
		Amount:    400,
		Type:      models.Credit,
	})
	if err != nil {
		t.Error(err)
	}

	debitTransactionJSON, err := json.Marshal(handlers.TransactinRequest{
		AccountID: accountID,
		Reference: uuid.Must(uuid.NewV4()),
		Amount:    401,
		Type:      models.Debit,
	})

	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(transactionJSON))
	router.ServeHTTP(sw, req)
	var response SuccessResponse
	err = json.Unmarshal(sw.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusCreated, sw.Code)

	ew := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(debitTransactionJSON))
	router.ServeHTTP(ew, req)
	var errorResponse ErrorResponse
	err = json.Unmarshal(ew.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, ew.Code)
	assert.Equal(t, errorResponse.Message, "Invalid transaction amount")

}

func TestValidationErrorOnTransaction(t *testing.T) {
	router := SetupRouter()

	transactionJSON, err := json.Marshal(handlers.TransactinRequest{
		AccountID: accountID,
		Reference: uuid.Must(uuid.NewV4()),
	})

	if err != nil {
		t.Error(err)
	}
	ew := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(transactionJSON))
	router.ServeHTTP(ew, req)
	var errorResponse ErrorResponse
	err = json.Unmarshal(ew.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusBadRequest, ew.Code)
	assert.Equal(t, errorResponse.Message, "Request validation failed")
	assert.Equal(t, errorResponse.Error, "Amount is required")

}

func TestValidSingleTransaction(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()

	transactionJSON, err := json.Marshal(handlers.TransactinRequest{
		AccountID: accountID,
		Reference: uuid.Must(uuid.NewV4()),
		Amount:    400,
		Type:      models.Credit,
	})

	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(transactionJSON))
	router.ServeHTTP(w, req)
	var response SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, response.Status, true)

	w = httptest.NewRecorder()
	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/v1/transactions/%s", response.Data.ID), nil)
	router.ServeHTTP(w, getReq)
	var singleTransaction SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &singleTransaction)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, singleTransaction.Message, "Retrieved transaction")
	assert.Equal(t, singleTransaction.Data.ID, response.Data.ID)
}
