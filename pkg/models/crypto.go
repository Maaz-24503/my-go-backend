package models

import "time"

// CoinGeckoResponse : CoinGecko API response structure
type CoinGeckoResponse struct {
	ID                    string  `json:"id"`
	Symbol                string  `json:"symbol"`
	Name                  string  `json:"name"`
	CurrentPrice          float64 `json:"current_price"`
	MarketCap             int64   `json:"market_cap"`
	MarketCapRank         int     `json:"market_cap_rank"`
	PriceChange24h        float64 `json:"price_change_24h"`
	PriceChangePercent24h float64 `json:"price_change_percentage_24h"`
	LastUpdated           string  `json:"last_updated"`
}

// CryptoData : Our internal crypto data structure
type CryptoData struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	MarketCap     int64     `json:"market_cap"`
	Rank          int       `json:"rank"`
	Change24h     float64   `json:"change_24h"`
	ChangePercent float64   `json:"change_percent_24h"`
	FetchedAt     time.Time `json:"fetched_at"`
	Error         string    `json:"error,omitempty"`
}

// PortfolioRequest : Portfolio request/response models
type PortfolioRequest struct {
	Coins []string `json:"coins" binding:"required"`
}

type PortfolioResponse struct {
	Portfolio    []CryptoData `json:"portfolio"`
	TotalValue   float64      `json:"total_value"`
	SuccessCount int          `json:"success_count"`
	ErrorCount   int          `json:"error_count"`
	FetchTime    string       `json:"fetch_time"`
}

// BulkCryptoRequest : Bulk crypto request
type BulkCryptoRequest struct {
	Coins   []string `json:"coins" binding:"required"`
	Timeout int      `json:"timeout,omitempty"` // seconds
}

type StreamEvent struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ID        string      `json:"id,omitempty"`
}

type PriceUpdate struct {
	CoinID     string    `json:"coin_id"`
	Symbol     string    `json:"symbol"`
	Price      float64   `json:"price"`
	Change24h  float64   `json:"change_24h"`
	Timestamp  time.Time `json:"timestamp"`
	UpdateType string    `json:"update_type"` // "price", "volume", "market_cap"
}

type StreamConfig struct {
	Coins      []string      `json:"coins"`
	Interval   time.Duration `json:"interval"`
	MaxUpdates int           `json:"max_updates,omitempty"`
}

type WebSocketMessage struct {
	Action string      `json:"action"` // "subscribe", "unsubscribe", "ping"
	Data   interface{} `json:"data"`
	ID     string      `json:"id,omitempty"`
}
