package main

import (
	"log"

	"github.com/CeoFred/fairmoney/constants"
	"github.com/CeoFred/fairmoney/database"
	"github.com/CeoFred/fairmoney/internal/helpers"
	"github.com/CeoFred/fairmoney/internal/routes"

	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "golang.org/x/text"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	g := gin.Default()
	v := constants.New()

	// Parse command-line flags
	flag.Parse()
	_ = helpers.NewCache()

	// Middleware
	g.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	g.Use(gin.Logger())

	g.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	g.MaxMultipartMemory = 8 << 20

	dbConfig := database.Config{
		Host:     v.DbHost,
		Port:     v.DbPort,
		Password: v.DbPassword,
		User:     v.DbUser,
		DBName:   v.DbName,
	}
	database.Connect(&dbConfig)
	database.RunManualMigration(database.DB)
	// Set up Swagger documentation

	v1 := g.Group("/v1")

	// Bind routes
	routes.Routes(v1, database.DB)

	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Whooops! Not Found",
		})
	})

	return g
}

func main() {
	constant := constants.New()

	g := SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = constant.Port
	}

	// Listen on port set in .env
	log.Fatal(g.Run(":" + port))
}
