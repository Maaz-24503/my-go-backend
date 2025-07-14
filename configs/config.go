package configs

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Host         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	JWTSecret    string
	JWTExpiresIn time.Duration
	AppEnv       string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	jwtExpires, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))

	return &Config{
		Port:         getEnv("PORT", "8095"),
		Host:         getEnv("HOST", "localhost"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "30532"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "myapp"),
		DBSSLMode:    getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:    getEnv("JWT_SECRET", "tHiSiSaSeCrEt"),
		JWTExpiresIn: jwtExpires,
		AppEnv:       getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
