package config

import "os"

var (
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	JWTSecret  string
	AWSRegion  string
	Port       string
)

func Init() {
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBName = os.Getenv("DB_NAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	JWTSecret = os.Getenv("JWT_SECRET")
	AWSRegion = os.Getenv("AWS_REGION")
	Port = os.Getenv("PORT")
}
