package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	BackendPort string
	MySQLHost   string
	MySQLPort   string
	MySQLUser   string
	MySQLPass   string
	MySQLDB     string
}

func Load(envFile string) *Config {

	if err := godotenv.Load(envFile); err != nil {
		if os.Getenv("APP_ENV") == "prod" {
			log.Println("[INFO] Running in production mode. Using system environment variables.")
		}
	}

	return &Config{
		AppEnv:      mustGetVal("APP_ENV"),
		BackendPort: mustGetVal("BACKEND_PORT"),
	}
}

func mustGetVal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[FATAL] Environment variable '%s' is missing", key)
	}
	return val
}
