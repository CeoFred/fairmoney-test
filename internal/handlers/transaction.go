package handlers

import (
	// "context"
	// "fmt"
	"log"
	"net/http"

	// "time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"github.com/CeoFred/fairmoney/acme_api_mock"
	"github.com/CeoFred/fairmoney/internal/helpers"
	"github.com/CeoFred/fairmoney/internal/models"
	"github.com/CeoFred/fairmoney/internal/repository"
)

var ()

type TransactionHandler struct {
	accountRepository     *repository.AccountRepository
	transactionRepository *repository.TransactionRepository
}

func NewTransactionHandler(accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
) *TransactionHandler {
	return &TransactionHandler{
		accountRepository:     accountRepo,
		transactionRepository: transactionRepo,
	}
}

type TransactinRequest struct {
	AccountID uuid.UUID                `json:"account_id" validate:"required"`
	Reference uuid.UUID                `json:"reference" validate:"required"`
	Amount    float64                  `json:"amount" validate:"required"`
	Type      models.TransactionIntent `json:"type" validate:"required"`
	ForceDuplicate bool `json:"force_duplicate"`
}

// SingleTransaction is a route handler that returns a single transaction using a unique reference ID
//
// @Router /transactions/{:transaction_id} [get]
func (u *TransactionHandler) SingleTransaction(c *gin.Context) {

	transaction_id := c.Param("transaction_id")

	if transaction_id == "" {
		helpers.ReturnJSON(c, "Transaction ID is required", nil, http.StatusBadRequest)
	}

	// fetch record from acme
	apiKey := "test-api_key"
	client := acme_http_mock.NewHTTPClient(&apiKey)

	responseBody, err := client.SendGetRequest("/v1/payments/", uuid.FromStringOrNil(transaction_id))
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusBadGateway)
		return
	}

	if responseBody == nil {
		helpers.ReturnError(c, "Transaction not found", nil, http.StatusNotFound)
		return
	}

	transaction, err := u.transactionRepository.Find(responseBody.Reference)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if transaction == nil {
		helpers.ReturnError(c, "Transaction not found", nil, http.StatusNotFound)
		return
	}

	helpers.ReturnJSON(c, "Retrieved transaction", transaction, http.StatusOK)
}

// NewTransaction is a route handler that creates a new account transaction
//
// @Router /transactions [post]
func (u *TransactionHandler) NewTransaction(c *gin.Context) {
	var input TransactinRequest
	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnJSON(c, "Failed to retrieve validated request body", nil, http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(TransactinRequest)
	if !ok {
		helpers.ReturnJSON(c, "Failed to convert types", nil, http.StatusBadRequest)
		return
	}
	// 

	account, err := u.accountRepository.Find(input.AccountID)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if account == nil {
		helpers.ReturnJSON(c, "Account not found", nil, http.StatusNotFound)
		return
	}

	// check for existing transaction with input reference
	txnExists, err := u.transactionRepository.Exists(input.Reference)
	if err != nil {
		if err.Error() != "record not found" {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
			return
		}
	}

	if txnExists {
		log.Println("duplicate transaction reference", input.Reference)
		helpers.ReturnJSON(c, "Duplicate transaction reference", input.Reference, http.StatusBadRequest)
		return
	}

	// create a new transaction with pending status
	transaction := &models.Transaction{
		Amount:    input.Amount,
		ID:        input.Reference,
		AccountID: input.AccountID,
		Status:    models.PendingTransaction,
		Type:      input.Type,
	}

	err = u.transactionRepository.Create(transaction)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	// send request to acme service
	apiKey := "test-api_key"
	client := acme_http_mock.NewHTTPClient(&apiKey)

	txn := acme_http_mock.Transaction{
		Reference: transaction.ID,
		Amount:    transaction.Amount,
		AccountID: transaction.AccountID,
	}
	response, err := client.SendPostRequest("/v1/payments", &txn)
	if err != nil {

		// mark transaction as pending
		transaction.Status = models.PendingTransaction

		if _, err := u.transactionRepository.Save(transaction); err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
			return
		}

		helpers.ReturnError(c, "Gateway error", err, http.StatusBadGateway)
		return
	}

	if response != nil {
		// update transaction
		transaction.Status = models.SuccessfulTransaction

		switch transaction.Type {
		case models.Credit:
			account.Balance += transaction.Amount
		case models.Debit:
			balance := account.Balance - transaction.Amount
			if balance < 0 {
				// mark as failed
				transaction.Status = models.FailedTransaction
				if _, err := u.transactionRepository.Save(transaction); err != nil {
					helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
					return
				}

				helpers.ReturnJSON(c, "Invalid transaction amount", nil, http.StatusBadRequest)
				return
			} else {
				account.Balance -= transaction.Amount
			}

		default:
			transaction.Status = models.FailedTransaction
			log.Println("invalid transaction type", transaction.Type)
			helpers.ReturnJSON(c, "Invalid transaction intent", nil, http.StatusBadRequest)
			return
		}

		if _, err := u.transactionRepository.Save(transaction); err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
			return
		}

		_, err = u.accountRepository.Save(account)

		if err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
			return
		}
	}

	helpers.ReturnJSON(c, "Transaction record saved successfully", transaction, http.StatusCreated)
}
