package acme_http_mock

import (
	"fmt"

	"github.com/CeoFred/fairmoney/internal/helpers"
	"github.com/CeoFred/fairmoney/validator"

	"github.com/gofrs/uuid"

	"errors"
	"time"
)

type HTTPClient struct {
	apiKey       *string
	CacheManager *helpers.Cache
}

func NewHTTPClient(apiKey *string) *HTTPClient {
	return &HTTPClient{
		apiKey:       apiKey,
		CacheManager: helpers.CacheManager,
	}
}

func (c *HTTPClient) SendPostRequest(url string, data *Transaction) (*Transaction, error) {

	err := validator.Validate(data)
	if err != nil {
		return nil, err
	}

	if data.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// 1 second delay
	time.Sleep(1 * time.Second)

	// store data in cache
	c.CacheManager.Set(data.Reference.String(), data, time.Hour*6000)

	return data, nil
}

func (c *HTTPClient) SendGetRequest(url string, id uuid.UUID) (*Transaction, error) {

	// find record from cache
	txn, exist := c.CacheManager.Get(id.String())

	fmt.Println(c.CacheManager)
	if !exist {
		return nil, errors.New("no record found")
	}

	responseData, ok := txn.(*Transaction)

	if !ok {
		return nil, errors.New("no transaction found,failed to parse format")
	}
	return responseData, nil
}
