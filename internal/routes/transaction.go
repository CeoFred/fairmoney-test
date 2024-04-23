package routes

import (
	"github.com/CeoFred/fairmoney/internal/handlers"
	"github.com/CeoFred/fairmoney/internal/repository"
	"github.com/CeoFred/fairmoney/internal/validators"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func RegisterTransactionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	r := router.Group("transactions")

	handler := handlers.NewTransactionHandler(repository.NewAccountRepository(db), repository.NewTransactionRepository(db))

	r.POST("", validators.ValidateNewTransaction, handler.NewTransaction)
	r.GET("/:transaction_id", handler.SingleTransaction)

}
