package acme_http_mock

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/CeoFred/fairmoney/internal/helpers"
)

func TestSendPostRequest(t *testing.T) {
	client := NewHTTPClient(nil)
	client.CacheManager = helpers.NewCache()

	testCases := []struct {
		testName      string
		Transaction   *Transaction
		expectedError error
	}{
		{
			testName: "Valid transaction",
			Transaction: &Transaction{
				Reference: uuid.Must(uuid.NewV4()),
				Amount:    200.2,
				AccountID: uuid.Must(uuid.NewV4()),
			},
			expectedError: nil,
		},
		{
			testName: "Transaction with negative amount",
			Transaction: &Transaction{
				Reference: uuid.Must(uuid.NewV4()),
				Amount:    -50.0,
				AccountID: uuid.Must(uuid.NewV4()),
			},
			expectedError: errors.New("amount must be greater than zero"),
		},
		{
			testName: "Transaction with zero amount",
			Transaction: &Transaction{
				Reference: uuid.Must(uuid.NewV4()),
				Amount:    0.0,
				AccountID: uuid.Must(uuid.NewV4()),
			},
			expectedError: errors.New("Amount is required"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			returnedTransaction, err := client.SendPostRequest("/v1/payments", tc.Transaction)
			if err != nil && tc.expectedError != nil {
				rErr := err.Error()
				expectedErr := tc.expectedError.Error()

				if rErr != expectedErr {
					t.Errorf("SendPostRequest returned an unexpected error: got %v, want %v", rErr, expectedErr)
				}
			}

			if returnedTransaction != nil && returnedTransaction.Reference != tc.Transaction.Reference {
				t.Error("SendPostRequest did not return the expected transaction")
			}

			if tc.expectedError == nil {
				Existingtxn, err := client.SendGetRequest("/v1/payments", tc.Transaction.Reference)

				if err != nil {
					t.Errorf("SendGetRequest returned an unexpected error: got %v", err)
				}

				if Existingtxn != nil && Existingtxn.Reference != tc.Transaction.Reference {
					t.Errorf("SendGetRequest returned the expected transaction. expected %v, got %v", tc.Transaction.Reference, Existingtxn.Reference)
				}
			}

		})
	}

}
