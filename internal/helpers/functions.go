package helpers

import (
	"fmt"
	"net/http"

	"time"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Characters to choose from
	randomLength = 10                                                               // Length of the random string
)

func ReturnJSON(c *gin.Context, message string, data interface{}, statusCode int) {
	c.Status(statusCode)
	c.JSON(statusCode, gin.H{
		"status":  statusCode <= 201,
		"message": message,
		"data":    data,
	})
}

func ReturnError(c *gin.Context, message string, err error, status int) {
	c.JSON(status, gin.H{
		"message": message,
		"error":   err.Error(),
		"status":  false,
	})
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letterBytes[b%byte(len(letterBytes))]
	}

	return string(bytes), nil
}

func Dispatch200OK(c *gin.Context, message string, data interface{}) {
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Dispatch201Created(c *gin.Context, message string, data interface{}) {
	c.Status(http.StatusCreated)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// 500 - internal server error
func Dispatch500Error(c *gin.Context, err error) {
	c.Status(http.StatusInternalServerError)
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
		"data":    nil,
	})
}

// 400 - bad request
func Dispatch400Error(c *gin.Context, msg string, err any) {
	c.Status(http.StatusBadRequest)
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": msg,
		"data":    err,
	})
}

// 404 - not found
func Dispatch404Error(c *gin.Context, msg string, err any) {
	c.Status(http.StatusNotFound)
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": msg,
		"data":    err,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func TimeNow(timezone string) (string, error) {

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	currentTime := time.Now().In(location)
	return currentTime.String(), nil
}

type AppError struct {
	message string
}

func (e AppError) Error() string {
	return e.message
}

func NewError(message string) *AppError {
	return &AppError{message: message}
}

func GetBaseURL(c *gin.Context) string {
	scheme := "http" // Default scheme
	isLocal := gin.Mode() == gin.DebugMode

	if isLocal {
		// Running in local development mode
		scheme = "http"
	} else {
		// Running in production or other mode
		scheme = "https"
	}

	// Get the host (domain) from the request
	host := c.Request.Host

	// Construct the base URL by combining the scheme and host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	return baseURL
}
