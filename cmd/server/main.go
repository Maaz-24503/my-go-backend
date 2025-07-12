package main

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"my-go-backend/configs"
	"my-go-backend/internal/handlers"
	"my-go-backend/internal/services"
	"my-go-backend/pkg/models"
	"time"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Connect to database
	db, err := connectDatabase(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database tables
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(db, config.JWTSecret, config.JWTExpiresIn)
	userService := services.NewUserService(db)
	cryptoService := services.NewCryptoService()

	// Start background price streaming for WebSocket subscribers
	ctx := context.Background()
	popularCoins := []string{"bitcoin", "ethereum", "bnb", "solana", "cardano"}
	go cryptoService.StartPriceStreaming(ctx, popularCoins, 5*time.Second)

	// Setup routes
	router := handlers.SetupRoutes(authService, userService, cryptoService, config.JWTSecret)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func connectDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		config.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
