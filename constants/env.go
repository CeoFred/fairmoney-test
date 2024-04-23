package constants

import (
	"log"
	"os"
	// "regexp"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	Env        string
	ProjectID  string
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file %s", err)
	}

}

func New() *Config {
	return &Config{
		DbHost:     getEnv("POSTGRES_HOST", ""),
		DbUser:     getEnv("POSTGRES_USER", ""),
		DbPassword: getEnv("POSTGRES_PASSWORD", ""),
		DbName:     getEnv("POSTGRES_NAME", ""),
		DbPort:     getEnv("POSTGRES_PORT", ""),
		Port:       getEnv("PORT", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
