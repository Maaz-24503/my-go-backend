package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"my-go-backend/pkg/models"
)

type CryptoService struct {
	client  *resty.Client
	baseURL string
	// Mutex for thread-safe operations
	mu sync.RWMutex
	// In-memory cache with timestamp
	cache map[string]models.CryptoData

	subscribers map[string]chan models.StreamEvent // WebSocket subscribers
	subMu       sync.RWMutex                       // Protect subscribers map
}

func NewCryptoService() *CryptoService {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	return &CryptoService{
		client:      client,
		baseURL:     "https://api.coingecko.com/api/v3",
		cache:       make(map[string]models.CryptoData),
		subscribers: make(map[string]chan models.StreamEvent),
	}
}

// GetSingleCrypto fetches data for a single cryptocurrency
func (s *CryptoService) GetSingleCrypto(coinID string) (*models.CryptoData, error) {
	// Check cache first (with read lock)
	s.mu.RLock()
	if cached, exists := s.cache[coinID]; exists {
		// Cache valid for 1 minute
		if time.Since(cached.FetchedAt) < time.Minute {
			s.mu.RUnlock()
			log.Printf("Cache hit for %s", coinID)
			return &cached, nil
		}
	}
	s.mu.RUnlock()

	// Make API call
	url := fmt.Sprintf("%s/coins/markets", s.baseURL)

	var response []models.CoinGeckoResponse
	resp, err := s.client.R().
		SetQueryParam("vs_currency", "usd").
		SetQueryParam("ids", coinID).
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode())
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("coin not found: %s", coinID)
	}

	// Convert to our internal structure
	crypto := models.CryptoData{
		ID:            response[0].ID,
		Symbol:        response[0].Symbol,
		Name:          response[0].Name,
		Price:         response[0].CurrentPrice,
		MarketCap:     response[0].MarketCap,
		Rank:          response[0].MarketCapRank,
		Change24h:     response[0].PriceChange24h,
		ChangePercent: response[0].PriceChangePercent24h,
		FetchedAt:     time.Now(),
	}

	// Update cache (with write lock)
	s.mu.Lock()
	s.cache[coinID] = crypto
	s.mu.Unlock()

	return &crypto, nil
}

// GetBulkCrypto demonstrates goroutines, wait groups, and locks
func (s *CryptoService) GetBulkCrypto(coins []string, timeout time.Duration) (*models.PortfolioResponse, error) {
	startTime := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channels for collecting results
	results := make(chan models.CryptoData, len(coins))
	errors := make(chan error, len(coins))

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Launch goroutines for each coin
	for _, coin := range coins {
		wg.Add(1)
		go func(coinID string) {
			defer wg.Done()

			// Create a channel for this specific request
			done := make(chan struct{})
			var crypto *models.CryptoData
			var err error

			// Launch the actual API call in another goroutine
			go func() {
				defer close(done)
				crypto, err = s.GetSingleCrypto(coinID)
			}()

			// Wait for either completion or context timeout
			select {
			case <-done:
				if err != nil {
					log.Printf("Error fetching %s: %v", coinID, err)
					// Send error data instead of nil
					results <- models.CryptoData{
						ID:        coinID,
						Error:     err.Error(),
						FetchedAt: time.Now(),
					}
				} else {
					results <- *crypto
				}
			case <-ctx.Done():
				log.Printf("Timeout fetching %s", coinID)
				results <- models.CryptoData{
					ID:        coinID,
					Error:     "timeout",
					FetchedAt: time.Now(),
				}
			}
		}(coin)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var portfolio []models.CryptoData
	var totalValue float64
	successCount := 0
	errorCount := 0

	// Use mutex for thread-safe aggregation
	var resultMu sync.Mutex

	// Collect all results
	for result := range results {
		resultMu.Lock()
		portfolio = append(portfolio, result)

		if result.Error == "" {
			totalValue += result.Price
			successCount++
		} else {
			errorCount++
		}
		resultMu.Unlock()
	}

	return &models.PortfolioResponse{
		Portfolio:    portfolio,
		TotalValue:   totalValue,
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		FetchTime:    fmt.Sprintf("%.2fs", time.Since(startTime).Seconds()),
	}, nil
}

// GetPortfolioRealtime demonstrates different concurrency patterns
func (s *CryptoService) GetPortfolioRealtime(coins []string) (*models.PortfolioResponse, error) {
	startTime := time.Now()

	// Buffered channel to prevent blocking
	results := make(chan models.CryptoData, len(coins))

	// Use sync.WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Launch limited number of goroutines (rate limiting)
	semaphore := make(chan struct{}, 5) // Max 5 concurrent requests

	for _, coin := range coins {
		wg.Add(1)
		go func(coinID string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			log.Printf("Fetching %s...", coinID)

			crypto, err := s.GetSingleCrypto(coinID)
			if err != nil {
				results <- models.CryptoData{
					ID:        coinID,
					Error:     err.Error(),
					FetchedAt: time.Now(),
				}
				return
			}

			results <- *crypto
		}(coin)
	}

	// Close results channel when done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Thread-safe result collection
	var portfolio []models.CryptoData
	var totalValue float64
	successCount := 0
	errorCount := 0

	// Collect results
	for result := range results {
		portfolio = append(portfolio, result)

		if result.Error == "" {
			totalValue += result.Price
			successCount++
		} else {
			errorCount++
		}
	}

	return &models.PortfolioResponse{
		Portfolio:    portfolio,
		TotalValue:   totalValue,
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		FetchTime:    fmt.Sprintf("%.2fs", time.Since(startTime).Seconds()),
	}, nil
}

// ClearCache demonstrates write locks
func (s *CryptoService) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache = make(map[string]models.CryptoData)
	log.Println("Cache cleared")
}

// GetCacheStats demonstrates read locks
func (s *CryptoService) GetCacheStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"cached_coins": len(s.cache),
		"cache_keys": func() []string {
			keys := make([]string, 0, len(s.cache))
			for k := range s.cache {
				keys = append(keys, k)
			}
			return keys
		}(),
	}
}

// StreamPriceUpdates - Server-Sent Events streaming
func (s *CryptoService) StreamPriceUpdates(ctx context.Context, config models.StreamConfig) <-chan models.StreamEvent {
	eventChan := make(chan models.StreamEvent, 100)

	go func() {
		defer close(eventChan)

		ticker := time.NewTicker(config.Interval)
		defer ticker.Stop()

		updateCount := 0

		for {
			select {
			case <-ctx.Done():
				log.Println("Stream context cancelled")
				return
			case <-ticker.C:
				// Fetch latest prices for all coins concurrently
				s.streamPriceUpdates(config.Coins, eventChan)

				updateCount++
				if config.MaxUpdates > 0 && updateCount >= config.MaxUpdates {
					log.Printf("Reached max updates limit: %d", config.MaxUpdates)
					return
				}
			}
		}
	}()

	return eventChan
}

// streamPriceUpdates - Helper to fetch and send price updates
func (s *CryptoService) streamPriceUpdates(coins []string, eventChan chan<- models.StreamEvent) {
	var wg sync.WaitGroup
	updateChan := make(chan models.PriceUpdate, len(coins))

	// Fetch prices concurrently
	for _, coin := range coins {
		wg.Add(1)
		go func(coinID string) {
			defer wg.Done()

			crypto, err := s.GetSingleCrypto(coinID)
			if err != nil {
				log.Printf("Error fetching %s: %v", coinID, err)
				return
			}

			// Simulate price fluctuation (in real app, this would be actual API data)
			priceChange := (rand.Float64() - 0.5) * 0.02 // ±1% change
			newPrice := crypto.Price * (1 + priceChange)

			update := models.PriceUpdate{
				CoinID:     crypto.ID,
				Symbol:     crypto.Symbol,
				Price:      newPrice,
				Change24h:  crypto.Change24h,
				Timestamp:  time.Now(),
				UpdateType: "price",
			}

			updateChan <- update
		}(coin)
	}

	// Close update channel when all goroutines complete
	go func() {
		wg.Wait()
		close(updateChan)
	}()

	// Send updates to event channel
	for update := range updateChan {
		event := models.StreamEvent{
			Type:      "price_update",
			Data:      update,
			Timestamp: time.Now(),
			ID:        uuid.New().String(),
		}

		select {
		case eventChan <- event:
		case <-time.After(100 * time.Millisecond):
			log.Println("Event channel full, dropping update")
		}
	}
}

// WebSocket subscriber management
func (s *CryptoService) AddSubscriber(id string) <-chan models.StreamEvent {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	eventChan := make(chan models.StreamEvent, 100)
	s.subscribers[id] = eventChan

	log.Printf("Added subscriber: %s", id)
	return eventChan
}

func (s *CryptoService) RemoveSubscriber(id string) {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	if eventChan, exists := s.subscribers[id]; exists {
		close(eventChan)
		delete(s.subscribers, id)
		log.Printf("Removed subscriber: %s", id)
	}
}

// Broadcast to all WebSocket subscribers
func (s *CryptoService) BroadcastToSubscribers(event models.StreamEvent) {
	s.subMu.RLock()
	defer s.subMu.RUnlock()

	for id, eventChan := range s.subscribers {
		select {
		case eventChan <- event:
		case <-time.After(100 * time.Millisecond):
			log.Printf("Subscriber %s channel full, dropping event", id)
		}
	}
}

// StartPriceStreaming - Background service for WebSocket broadcasting
func (s *CryptoService) StartPriceStreaming(ctx context.Context, coins []string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.subMu.RLock()
			hasSubscribers := len(s.subscribers) > 0
			s.subMu.RUnlock()

			if !hasSubscribers {
				continue
			}

			// Fetch and broadcast updates
			go func() {
				var wg sync.WaitGroup
				for _, coin := range coins {
					wg.Add(1)
					go func(coinID string) {
						defer wg.Done()

						crypto, err := s.GetSingleCrypto(coinID)
						if err != nil {
							return
						}

						// Simulate price change
						priceChange := (rand.Float64() - 0.5) * 0.02
						newPrice := crypto.Price * (1 + priceChange)

						update := models.PriceUpdate{
							CoinID:     crypto.ID,
							Symbol:     crypto.Symbol,
							Price:      newPrice,
							Change24h:  crypto.Change24h,
							Timestamp:  time.Now(),
							UpdateType: "price",
						}

						event := models.StreamEvent{
							Type:      "price_update",
							Data:      update,
							Timestamp: time.Now(),
							ID:        uuid.New().String(),
						}

						s.BroadcastToSubscribers(event)
					}(coin)
				}
				wg.Wait()
			}()
		}
	}
}
