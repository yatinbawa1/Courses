package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBHost         string
	DBPort         string
	DBUser         string
	DBName         string
	DBPassword     string
	JWTSecret      string
	AWSRegion      string
	Port           string
	EMAIL_API_KEY  string
	REDIS_ADDR     string
	SECURE_COOKIES bool
)

func Init(l *log.Logger) {
	err := godotenv.Load()

	if err != nil {
		l.Printf("Unable to locate ENV File")
	}

	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBName = os.Getenv("DB_NAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	JWTSecret = os.Getenv("JWT_SECRET")
	AWSRegion = os.Getenv("AWS_REGION")
	Port = os.Getenv("PORT")
	EMAIL_API_KEY = os.Getenv("EMAIL_API")
	REDIS_ADDR = os.Getenv("REDIS_ADDR")
	SECURE_COOKIES = false
}
