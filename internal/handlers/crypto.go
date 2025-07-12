package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"my-go-backend/internal/services"
	"my-go-backend/pkg/models"
)

type CryptoHandler struct {
	cryptoService *services.CryptoService
	upgrader      websocket.Upgrader // WebSocket upgrader
}

func NewCryptoHandler(cryptoService *services.CryptoService) *CryptoHandler {
	return &CryptoHandler{
		cryptoService: cryptoService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo
			},
		},
	}
}

// GetSingleCrypto - Get data for a single cryptocurrency
func (h *CryptoHandler) GetSingleCrypto(c *gin.Context) {
	coinID := c.Param("coinId")
	if coinID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Coin ID is required",
		})
		return
	}

	crypto, err := h.cryptoService.GetSingleCrypto(coinID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch crypto data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Crypto data retrieved successfully",
		Data:    crypto,
	})
}

// GetBulkCrypto - Demonstrates goroutines with timeout
func (h *CryptoHandler) GetBulkCrypto(c *gin.Context) {
	var req models.BulkCryptoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if len(req.Coins) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "At least one coin is required",
		})
		return
	}

	if len(req.Coins) > 20 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Maximum 20 coins allowed",
		})
		return
	}

	// Default timeout of 15 seconds
	timeout := 15 * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	portfolio, err := h.cryptoService.GetBulkCrypto(req.Coins, timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch bulk crypto data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bulk crypto data retrieved successfully",
		Data:    portfolio,
	})
}

// GetPortfolioRealtime - Demonstrates goroutines with rate limiting
func (h *CryptoHandler) GetPortfolioRealtime(c *gin.Context) {
	var req models.PortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if len(req.Coins) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "At least one coin is required",
		})
		return
	}

	portfolio, err := h.cryptoService.GetPortfolioRealtime(req.Coins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch portfolio data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Portfolio data retrieved successfully",
		Data:    portfolio,
	})
}

// GetCacheStats - Demonstrates read locks
func (h *CryptoHandler) GetCacheStats(c *gin.Context) {
	stats := h.cryptoService.GetCacheStats()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Cache statistics retrieved",
		Data:    stats,
	})
}

// ClearCache - Demonstrates write locks
func (h *CryptoHandler) ClearCache(c *gin.Context) {
	h.cryptoService.ClearCache()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Cache cleared successfully",
	})
}

// GetPopularCoins - Get top cryptocurrencies (query params demo)
func (h *CryptoHandler) GetPopularCoins(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	// Popular coins to demonstrate bulk fetching
	popularCoins := []string{
		"bitcoin", "ethereum", "tether", "bnb", "solana",
		"usdc", "xrp", "dogecoin", "cardano", "avalanche-2",
		"chainlink", "polygon", "litecoin", "uniswap", "ethereum-classic",
	}

	// Limit the coins based on the request
	if limit < len(popularCoins) {
		popularCoins = popularCoins[:limit]
	}

	portfolio, err := h.cryptoService.GetPortfolioRealtime(popularCoins)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch popular coins",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Popular coins retrieved successfully",
		Data:    portfolio,
	})
}

// StreamPrices - Server-Sent Events endpoint
func (h *CryptoHandler) StreamPrices(c *gin.Context) {
	// Parse query parameters
	coinsParam := c.Query("coins")
	if coinsParam == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "coins parameter is required",
		})
		return
	}

	coins := strings.Split(coinsParam, ",")
	intervalStr := c.DefaultQuery("interval", "5") // Default 5 seconds
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval < 1 {
		interval = 5
	}

	maxUpdatesStr := c.DefaultQuery("max_updates", "0")
	maxUpdates, _ := strconv.Atoi(maxUpdatesStr)

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// Create context with cancellation
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// Configure streaming
	config := models.StreamConfig{
		Coins:      coins,
		Interval:   time.Duration(interval) * time.Second,
		MaxUpdates: maxUpdates,
	}

	// Start streaming
	eventChan := h.cryptoService.StreamPriceUpdates(ctx, config)

	// Write events to client
	for event := range eventChan {
		eventData, err := json.Marshal(event)
		if err != nil {
			log.Printf("Error marshaling event: %v", err)
			continue
		}

		// Proper SSE format with ID and event type
		fmt.Fprintf(c.Writer, "id: %s\n", event.ID)
		fmt.Fprintf(c.Writer, "event: %s\n", event.Type)
		fmt.Fprintf(c.Writer, "data: %s\n\n", eventData)
		c.Writer.Flush()

		// Check if client disconnected
		if ctx.Err() != nil {
			break
		}
	}
}

// WebSocketHandler - WebSocket endpoint
func (h *CryptoHandler) WebSocketHandler(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Generate unique subscriber ID
	subscriberID := uuid.New().String()
	log.Printf("New WebSocket connection: %s", subscriberID)

	// Add subscriber
	eventChan := h.cryptoService.AddSubscriber(subscriberID)
	defer h.cryptoService.RemoveSubscriber(subscriberID)

	// Handle client messages in separate goroutine
	go func() {
		defer func() {
			log.Printf("WebSocket read goroutine ended for %s", subscriberID)
		}()

		for {
			var msg models.WebSocketMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("WebSocket read error for %s: %v", subscriberID, err)
				break
			}

			log.Printf("Received message from %s: %+v", subscriberID, msg)

			switch msg.Action {
			case "ping":
				err := conn.WriteJSON(models.WebSocketMessage{
					Action: "pong",
					Data:   "Server is alive",
					ID:     msg.ID,
				})
				if err != nil {
					log.Printf("Error sending pong: %v", err)
				}
			case "subscribe":
				err := conn.WriteJSON(models.WebSocketMessage{
					Action: "subscribed",
					Data:   fmt.Sprintf("Subscribed to updates for %s", subscriberID),
					ID:     msg.ID,
				})
				if err != nil {
					log.Printf("Error sending subscribe confirmation: %v", err)
				}
			}
		}
	}()

	// Send events to client
	for event := range eventChan {
		err := conn.WriteJSON(event)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

// StreamPortfolio - Stream portfolio updates
func (h *CryptoHandler) StreamPortfolio(c *gin.Context) {
	var req models.PortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// Stream portfolio updates every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			portfolio, err := h.cryptoService.GetPortfolioRealtime(req.Coins)
			if err != nil {
				log.Printf("Error getting portfolio: %v", err)
				continue
			}

			event := models.StreamEvent{
				Type:      "portfolio_update",
				Data:      portfolio,
				Timestamp: time.Now(),
				ID:        uuid.New().String(),
			}

			eventData, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error marshaling portfolio event: %v", err)
				continue
			}

			fmt.Fprintf(c.Writer, "data: %s\n\n", eventData)
			c.Writer.Flush()
		}
	}
}
