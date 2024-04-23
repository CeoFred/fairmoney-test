package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.RouterGroup, db *gorm.DB) {

	RegisterTransactionRoutes(r, db)

}
