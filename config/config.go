package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, falling back to environment variables")
	}

	DBHost = os.Getenv("DB_HOST")
	if DBHost == "" {
		DBHost = "localhost"
	}
	DBPort = os.Getenv("DB_PORT")
	if DBPort == "" {
		DBPort = "5432"
	}
	DBUser = os.Getenv("DB_USER")
	if DBUser == "" {
		DBUser = "postgres"
	}
	DBPassword = os.Getenv("DB_PASSWORD")
	if DBPassword == "" {
		log.Fatal("DB_PASSWORD is required")
	}
	DBName = os.Getenv("DB_NAME")
	if DBName == "" {
		DBName = "polyclinic"
	}
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
}
