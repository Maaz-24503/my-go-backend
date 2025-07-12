package handlers

import (
	"github.com/gin-gonic/gin"
	"my-go-backend/internal/middleware"
	"my-go-backend/internal/services"
)

func SetupRoutes(
	authService *services.AuthService,
	userService *services.UserService,
	cryptoService *services.CryptoService,
	jwtSecret string,
) *gin.Engine {
	router := gin.Default()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check (no auth required)
	router.GET("/health", HealthCheck)

	// API v1 group
	v1 := router.Group("/api/v1")

	// Auth routes (no auth required)
	authHandler := NewAuthHandler(authService)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// User routes (auth required)
	userHandler := NewUserHandler(userService)
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtSecret))
	{
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	cryptoHandler := NewCryptoHandler(cryptoService)
	crypto := v1.Group("/crypto")
	crypto.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Single crypto data
		crypto.GET("/:coinId", cryptoHandler.GetSingleCrypto)

		// Bulk operations (demonstrates goroutines)
		crypto.POST("/bulk", cryptoHandler.GetBulkCrypto)
		crypto.POST("/portfolio", cryptoHandler.GetPortfolioRealtime)

		// Popular coins (query params)
		crypto.GET("/popular", cryptoHandler.GetPopularCoins)

		// Cache operations (demonstrates locks)
		crypto.GET("/cache/stats", cryptoHandler.GetCacheStats)
		crypto.DELETE("/cache", cryptoHandler.ClearCache)

		// Streaming routes
		crypto.GET("/stream/prices", cryptoHandler.StreamPrices)        // SSE
		crypto.GET("/stream/ws", cryptoHandler.WebSocketHandler)        // WebSocket
		crypto.POST("/stream/portfolio", cryptoHandler.StreamPortfolio) // JSON streaming
	}

	return router
}
