package validators

import (
	"net/http"

	"github.com/CeoFred/fairmoney/internal/handlers"
	"github.com/CeoFred/fairmoney/internal/helpers"
	"github.com/CeoFred/fairmoney/validator"

	"github.com/gin-gonic/gin"
)

func ValidateNewTransaction(c *gin.Context) {
	var body handlers.TransactinRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		helpers.ReturnError(c, "Invalid Payload", err, http.StatusBadRequest)
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.ReturnError(c, "Request validation failed", err, http.StatusBadRequest)
		c.Abort()
		return
	}
	c.Set("validatedRequestBody", body)
	c.Next()
}
