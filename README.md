# Go Backend Learning Project: Crypto Portfolio Tracker

![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

A comprehensive Go backend application demonstrating modern Go development practices, concurrency patterns, and real-time data streaming. This project serves as a practical learning resource for Go backend development with a crypto portfolio tracking system.

## 🎯 Learning Objectives

This project teaches you:
- **Go Project Structure**: Industry-standard organization patterns
- **Concurrency**: Goroutines, channels, wait groups, and mutex
- **REST API Development**: Using Gin framework
- **Real-time Streaming**: Server-Sent Events (SSE) and WebSockets
- **Database Integration**: PostgreSQL with GORM
- **Authentication**: JWT-based auth system
- **Clean Architecture**: Separation of concerns with layers
- **Error Handling**: Proper Go error handling patterns
- **Configuration Management**: Environment-based configuration
- **Caching**: In-memory caching with thread safety
- **Rate Limiting**: Controlling external API calls

## 🏗️ Project Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

```
my-go-backend/
├── cmd/                    # Application entry points
├── configs/                # Configuration management
├── internal/               # Private application code
│   ├── handlers/          # HTTP request handlers (controllers)
│   ├── middleware/        # HTTP middleware
│   └── services/          # Business logic layer
├── pkg/                   # Public library code
│   └── models/           # Data structures and models
├── docs/                  # Documentation
├── .env                   # Environment variables
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md             # This file
```

### 🔍 Directory Purposes Explained

#### `cmd/` - Application Entry Points
Contains the main applications for this project. The directory name for each application should match the name of the executable you want to have.

**Why separate?** Large projects often have multiple executables (web server, CLI tools, workers). Each gets its own subdirectory.

```go
// cmd/server/main.go - Web server entry point
func main() {
    // Initialize configuration
    // Connect to database  
    // Start HTTP server
}
```

#### `internal/` - Private Application Code
Code that you don't want others importing. Go compiler enforces this - other projects cannot import from internal directories.

**Key Concept**: Everything in `internal/` is implementation detail, not public API.

#### `pkg/` - Public Library Code  
Code that's ok for other applications to import. Contains reusable components and shared models.

**Rule of Thumb**: If another project might use this code, put it in `pkg/`. If it's specific to this application, put it in `internal/`.

#### `configs/` - Configuration Management
Centralized configuration handling. Loads from environment variables, config files, or command-line flags.

## 🚀 Quick Start Guide

### Prerequisites

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **PostgreSQL**: Running on port 30532 (Docker/Kubernetes)
- **Git**: For version control

### Installation Steps

1. **Clone the repository**
```bash
git clone https://github.com/Maaz-24503/my-go-backend.git
cd my-go-backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Run the application**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 🔧 Environment Configuration

Create a `.env` file in the root directory:

```env
# Server Configuration
PORT=8080
HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database_name
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h

# Application Environment
APP_ENV=development
```

**Learning Point**: Never commit `.env` files to version control. Use `.env.example` as a template.

## 📚 Detailed Architecture Guide

### Layer Responsibilities

```
┌─────────────────┐
│   HTTP Layer    │ ← Handlers (Gin routes, middleware)
├─────────────────┤
│  Business Layer │ ← Services (business logic, external APIs)  
├─────────────────┤
│    Data Layer   │ ← Models (structs, database schemas)
├─────────────────┤
│ Infrastructure  │ ← Database, configuration, external services
└─────────────────┘
```

### 🎯 Understanding Go Packages

#### What is a Package?
A package is a collection of Go source files in the same directory that are compiled together. 

```go
// All files in the same directory must have the same package declaration
package handlers  // This is the package name
```

#### Import Paths
```go
import (
    "fmt"                                    // Standard library
    "github.com/gin-gonic/gin"              // External package
    "my-go-backend/internal/services"       // Internal package
    "my-go-backend/pkg/models"              // Public package
)
```

**Learning Rule**: The import path after your module name must match the directory structure.

## 🔍 File-by-File Breakdown

Let me break down each file type and its purpose:

### Configuration Files

#### `go.mod` - Module Definition
```go
module my-go-backend

go 1.24.4

require (
    github.com/gin-gonic/gin v1.10.1
    github.com/golang-jwt/jwt/v5 v5.2.2
    // ... other dependencies
)
```

**Purpose**: 
- Defines your module name
- Specifies Go version requirement
- Lists all dependencies with their versions

**Learning Point**: Go modules replaced the old GOPATH system. The module name is used as the base for all import paths.

#### `go.sum` - Dependency Checksums
Contains cryptographic checksums of dependencies to ensure reproducible builds.

**Never edit manually** - Go manages this file automatically.

### Application Entry Point

#### `cmd/server/main.go` - Main Application
```go
func main() {
    // 1. Load configuration from environment
    config := configs.LoadConfig()
    
    // 2. Connect to database
    db, err := connectDatabase(config)
    
    // 3. Initialize services (business logic layer)
    authService := services.NewAuthService(db, config.JWTSecret, config.JWTExpiresIn)
    
    // 4. Setup HTTP routes and middleware
    router := handlers.SetupRoutes(authService, userService, cryptoService, config.JWTSecret)
    
    // 5. Start HTTP server
    router.Run(serverAddr)
}
```

**Purpose**: 
- Bootstraps the entire application
- Connects all layers together
- Handles graceful startup and shutdown

**Learning Pattern**: Main should be thin - it orchestrates but doesn't contain business logic.

### Configuration Layer

#### `configs/config.go` - Configuration Management
```go
type Config struct {
    Port         string
    Host         string
    DBHost       string
    DBPort       string
    // ... other configuration fields
}

func LoadConfig() *Config {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Return configuration with environment variables or defaults
    return &Config{
        Port: getEnv("PORT", "8080"),
        Host: getEnv("HOST", "localhost"),
        // ...
    }
}
```

**Key Concepts**:
- **Struct-based configuration**: Type-safe configuration management
- **Environment variable fallbacks**: Use env vars or sensible defaults
- **Single source of truth**: All config in one place

**Learning Point**: The `getEnv()` helper function provides default values when environment variables aren't set.

```go
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### Data Models Layer

#### `pkg/models/user.go` - User Data Structures
```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

**Understanding Struct Tags**:
- `json:"id"`: How field appears in JSON responses
- `json:"-"`: Exclude from JSON (security for passwords)
- `gorm:"primaryKey"`: Database primary key
- `gorm:"unique;not null"`: Database constraints

**Request/Response Models**:
```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

**Learning Pattern**: Separate request/response models from database models for:
- **Security**: Don't expose internal fields
- **Validation**: Different validation rules for input vs output
- **API Stability**: Change database without breaking API

#### `pkg/models/crypto.go` - Cryptocurrency Models
```go
type CryptoData struct {
    ID            string    `json:"id"`
    Symbol        string    `json:"symbol"`
    Name          string    `json:"name"`
    Price         float64   `json:"price"`
    MarketCap     int64     `json:"market_cap"`
    FetchedAt     time.Time `json:"fetched_at"`
    Error         string    `json:"error,omitempty"`
}
```

**Key Design Decisions**:
- `Error` field for handling individual failures in bulk operations
- `FetchedAt` for cache invalidation
- `omitempty` tag excludes empty errors from JSON

#### `pkg/models/response.go` - API Response Patterns
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

**Learning Pattern**: Consistent API response structure:
- Always include `success` boolean
- Descriptive `message` for humans
- Optional `data` for successful responses
- Optional `error` for error details

### Middleware Layer

#### `internal/middleware/auth.go` - JWT Authentication
```go
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract token from Authorization header
        authHeader := c.GetHeader("Authorization")
        
        // 2. Validate Bearer token format
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }
        
        // 3. Parse and validate JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        
        // 4. Extract claims and store in context
        claims, ok := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        
        // 5. Continue to next handler
        c.Next()
    }
}
```

**Middleware Concepts**:
- **Chain of Responsibility**: Each middleware can modify request/response
- `c.Next()`: Continue to next middleware/handler
- `c.Abort()`: Stop the chain and return immediately
- `c.Set()`: Store data for later handlers to use

#### `internal/middleware/cors.go` - Cross-Origin Resource Sharing
```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

**CORS Explanation**:
- Browsers block cross-origin requests by default
- CORS headers tell browser which origins are allowed
- OPTIONS requests are "preflight" checks

#### `internal/middleware/logger.go` - Request Logging
```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next() // Process request
        
        latency := time.Since(start)
        log.Printf("[%s] %s %s - %d - %v",
            c.Request.Method,
            c.Request.URL.Path,
            c.ClientIP(),
            c.Writer.Status(),
            latency,
        )
    }
}
```

**Logging Pattern**:
- Capture start time before processing
- Call `c.Next()` to process request
- Calculate metrics after processing
- Log useful information for debugging

## 🧠 Understanding Go Concurrency

This project demonstrates key Go concurrency concepts:

### Goroutines - Lightweight Threads
```go
// Launch goroutine for each coin
for _, coin := range coins {
    wg.Add(1)
    go func(coinID string) {
        defer wg.Done()
        
        // Fetch crypto data concurrently
        crypto, err := s.GetSingleCrypto(coinID)
        results <- crypto
    }(coin)
}
```

**Key Points**:
- Goroutines are much lighter than OS threads
- Pass variables as parameters to avoid closure issues
- Always capture loop variables in goroutine parameters

### WaitGroups - Synchronization
```go
var wg sync.WaitGroup

// Add to wait group before launching goroutine
wg.Add(1)

go func() {
    defer wg.Done() // Signal completion
    // Do work...
}()

// Wait for all goroutines to complete
wg.Wait()
```

**WaitGroup Methods**:
- `Add(delta)`: Increase counter
- `Done()`: Decrease counter by 1 (usually in defer)
- `Wait()`: Block until counter reaches 0

### Channels - Communication
```go
// Buffered channel prevents blocking
results := make(chan models.CryptoData, len(coins))

// Send to channel
results <- cryptoData

// Receive from channel
data := <-results

// Range over channel (until closed)
for data := range results {
    // Process data
}

// Close channel to signal no more data
close(results)
```

**Channel Patterns**:
- **Unbuffered**: Synchronous communication (sender blocks until receiver ready)
- **Buffered**: Asynchronous up to buffer size
- **Closing**: Signals no more data will be sent

### Mutexes - Thread Safety
```go
type CryptoService struct {
    mu    sync.RWMutex
    cache map[string]models.CryptoData
}

// Read operation (multiple readers allowed)
func (s *CryptoService) GetCacheStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return map[string]interface{}{
        "cached_coins": len(s.cache),
    }
}

// Write operation (exclusive access)
func (s *CryptoService) ClearCache() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.cache = make(map[string]models.CryptoData)
}
```

**Mutex Types**:
- `sync.Mutex`: Exclusive access (one reader OR one writer)
- `sync.RWMutex`: Multiple readers OR one writer (not both)

### Context - Cancellation and Timeouts
```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Use context in operations
select {
case result := <-workChan:
    // Work completed
case <-ctx.Done():
    // Timeout or cancellation
    return ctx.Err()
}
```

**Context Usage**:
- Pass context as first parameter to functions
- Check `ctx.Done()` for cancellation
- Always call `cancel()` to free resources

## 🏢 Business Logic Layer

### Services - The Heart of Your Application

Services contain the core business logic and coordinate between different layers.

#### `internal/services/auth.go` - Authentication Service
```go
type AuthService struct {
    db        *gorm.DB
    jwtSecret string
    jwtExpiry time.Duration
}

func NewAuthService(db *gorm.DB, jwtSecret string, jwtExpiry time.Duration) *AuthService {
    return &AuthService{
        db:        db,
        jwtSecret: jwtSecret,
        jwtExpiry: jwtExpiry,
    }
}

func (s *AuthService) Register(req *models.CreateUserRequest) (*models.UserResponse, error) {
    // 1. Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // 2. Create user in database
    user := models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
    }

    if err := s.db.Create(&user).Error; err != nil {
        return nil, err
    }

    // 3. Return safe user response (no password)
    return &models.UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
    }, nil
}
```

**Service Design Patterns**:
- **Constructor functions**: `NewAuthService()` initializes dependencies
- **Method receivers**: `(s *AuthService)` makes methods belong to the service
- **Error handling**: Always return errors as the last return value
- **Data transformation**: Convert between internal and external models

#### `internal/services/crypto.go` - Cryptocurrency Service
```go
type CryptoService struct {
    client      *resty.Client
    baseURL     string
    mu          sync.RWMutex
    cache       map[string]models.CryptoData
    subscribers map[string]chan models.StreamEvent
    subMu       sync.RWMutex
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
```

**Advanced Patterns Demonstrated**:

##### 1. Caching with Thread Safety
```go
func (s *CryptoService) GetSingleCrypto(coinID string) (*models.CryptoData, error) {
    // Check cache first (read lock)
    s.mu.RLock()
    if cached, exists := s.cache[coinID]; exists {
        if time.Since(cached.FetchedAt) < time.Minute {
            s.mu.RUnlock()
            return &cached, nil
        }
    }
    s.mu.RUnlock()

    // Fetch from API if not cached
    // ... API call logic ...

    // Update cache (write lock)
    s.mu.Lock()
    s.cache[coinID] = crypto
    s.mu.Unlock()

    return &crypto, nil
}
```

**Caching Strategy**:
- **Read-heavy optimization**: Use `RWMutex` for better read performance
- **TTL (Time To Live)**: Cache expires after 1 minute
- **Cache-aside pattern**: Check cache first, then fallback to API

##### 2. Concurrent API Calls with Error Handling
```go
func (s *CryptoService) GetBulkCrypto(coins []string, timeout time.Duration) (*models.PortfolioResponse, error) {
    startTime := time.Now()
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // Results channel
    results := make(chan models.CryptoData, len(coins))
    var wg sync.WaitGroup
    
    // Launch goroutines for each coin
    for _, coin := range coins {
        wg.Add(1)
        go func(coinID string) {
            defer wg.Done()
            
            done := make(chan struct{})
            var crypto *models.CryptoData
            var err error
            
            // Actual API call in separate goroutine
            go func() {
                defer close(done)
                crypto, err = s.GetSingleCrypto(coinID)
            }()
            
            // Wait for completion or timeout
            select {
            case <-done:
                if err != nil {
                    results <- models.CryptoData{
                        ID:    coinID,
                        Error: err.Error(),
                    }
                } else {
                    results <- *crypto
                }
            case <-ctx.Done():
                results <- models.CryptoData{
                    ID:    coinID,
                    Error: "timeout",
                }
            }
        }(coin)
    }
    
    // Close results channel when all goroutines complete
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results safely
    var portfolio []models.CryptoData
    var totalValue float64
    successCount := 0
    errorCount := 0
    
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
```

**Concurrency Patterns Explained**:
- **Fan-out**: Launch multiple goroutines from one place
- **Fan-in**: Collect results from multiple goroutines
- **Timeout handling**: Use context to cancel long-running operations
- **Error aggregation**: Collect both successful and failed results

##### 3. Rate Limiting with Semaphores
```go
func (s *CryptoService) GetPortfolioRealtime(coins []string) (*models.PortfolioResponse, error) {
    results := make(chan models.CryptoData, len(coins))
    var wg sync.WaitGroup
    
    // Semaphore for rate limiting
    semaphore := make(chan struct{}, 5) // Max 5 concurrent requests
    
    for _, coin := range coins {
        wg.Add(1)
        go func(coinID string) {
            defer wg.Done()
            
            // Acquire semaphore
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // Do work...
            crypto, err := s.GetSingleCrypto(coinID)
            // ... handle result ...
        }(coin)
    }
    
    // ... rest of implementation
}
```

**Rate Limiting Concept**:
- **Semaphore pattern**: Use buffered channel as counting semaphore
- **Resource protection**: Prevent overwhelming external APIs
- **Graceful degradation**: System remains responsive under load

## 🌐 HTTP Handler Layer

### Request Processing Pipeline

```
HTTP Request → Middleware Chain → Handler → Service → Database
                     ↓
HTTP Response ← JSON Serialization ← Business Logic ← Data Access
```

#### `internal/handlers/crypto.go` - HTTP Request Handlers
```go
type CryptoHandler struct {
    cryptoService *services.CryptoService
    upgrader      websocket.Upgrader
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
```

##### Basic REST Handler Pattern
```go
func (h *CryptoHandler) GetSingleCrypto(c *gin.Context) {
    // 1. Extract and validate parameters
    coinID := c.Param("coinId")
    if coinID == "" {
        c.JSON(http.StatusBadRequest, models.APIResponse{
            Success: false,
            Message: "Coin ID is required",
        })
        return
    }

    // 2. Call business logic
    crypto, err := h.cryptoService.GetSingleCrypto(coinID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.APIResponse{
            Success: false,
            Message: "Failed to fetch crypto data",
            Error:   err.Error(),
        })
        return
    }

    // 3. Return successful response
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Crypto data retrieved successfully",
        Data:    crypto,
    })
}
```

**Handler Responsibilities**:
- **Input validation**: Check required parameters
- **Request binding**: Parse JSON/form data into structs
- **Service coordination**: Call appropriate business logic
- **Response formatting**: Convert to consistent API response format
- **Error handling**: Return appropriate HTTP status codes

##### Advanced Handler: Request Validation
```go
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

    // Business rule validation
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

    // ... rest of handler logic
}
```

**Validation Layers**:
- **Struct validation**: Using `binding` tags
- **Business rules**: Custom validation logic
- **Security limits**: Prevent abuse (max 20 coins)

## 🔄 Real-time Streaming Implementation

This project demonstrates three streaming approaches:

### 1. Server-Sent Events (SSE)
```go
func (h *CryptoHandler) StreamPrices(c *gin.Context) {
    // Set SSE headers
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("Access-Control-Allow-Origin", "*")

    // Parse query parameters
    coinsParam := c.Query("coins")
    coins := strings.Split(coinsParam, ",")
    
    // Create streaming configuration
    config := models.StreamConfig{
        Coins:      coins,
        Interval:   time.Duration(interval) * time.Second,
        MaxUpdates: maxUpdates,
    }

    // Create context for cancellation
    ctx, cancel := context.WithCancel(c.Request.Context())
    defer cancel()

    // Start streaming
    eventChan := h.cryptoService.StreamPriceUpdates(ctx, config)

    // Stream events to client
    for event := range eventChan {
        eventData, _ := json.Marshal(event)
        
        // SSE format
        fmt.Fprintf(c.Writer, "id: %s\n", event.ID)
        fmt.Fprintf(c.Writer, "event: %s\n", event.Type)
        fmt.Fprintf(c.Writer, "data: %s\n\n", eventData)
        
        c.Writer.Flush()
        
        if ctx.Err() != nil {
            break
        }
    }
}
```

**SSE Characteristics**:
- **One-way communication**: Server to client only
- **HTTP-based**: Works through firewalls and proxies
- **Automatic reconnection**: Browser handles reconnection
- **Simple protocol**: Text-based event format

### 2. WebSocket Streaming
```go
func (h *CryptoHandler) WebSocketHandler(c *gin.Context) {
    // Upgrade HTTP connection to WebSocket
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer conn.Close()

    // Generate unique client ID
    clientID := uuid.New().String()
    
    // Add subscriber
    eventChan := h.cryptoService.AddSubscriber(clientID)
    defer h.cryptoService.RemoveSubscriber(clientID)

    // Handle incoming messages from client
    go func() {
        for {
            var msg models.WebSocketMessage
            err := conn.ReadJSON(&msg)
            if err != nil {
                log.Printf("WebSocket read error: %v", err)
                return
            }

            switch msg.Action {
            case "ping":
                conn.WriteJSON(models.WebSocketMessage{
                    Action: "pong",
                    Data:   "pong",
                    ID:     msg.ID,
                })
            case "subscribe":
                // Handle subscription logic
            }
        }
    }()

    // Handle outgoing messages to client
    for event := range eventChan {
        if err := conn.WriteJSON(event); err != nil {
            log.Printf("WebSocket write error: %v", err)
            return
        }
    }
}
```

**WebSocket Characteristics**:
- **Full-duplex communication**: Bidirectional real-time communication
- **Low overhead**: Binary protocol, minimal headers
- **Connection state**: Persistent connection with state management
- **Custom protocols**: Can implement custom message formats

### 3. Subscriber Management for WebSockets
```go
// In CryptoService
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
```

**Subscriber Pattern Benefits**:
- **Dynamic subscription**: Clients can connect/disconnect anytime
- **Resource cleanup**: Properly close channels and remove subscribers
- **Non-blocking broadcast**: Don't block if subscriber is slow

## 🛣️ Routing and Middleware Stack

### `internal/handlers/routes.go` - Route Configuration
```go
func SetupRoutes(
    authService *services.AuthService,
    userService *services.UserService,
    cryptoService *services.CryptoService,
    jwtSecret string,
) *gin.Engine {
    router := gin.Default()

    // Global middleware (applied to all routes)
    router.Use(middleware.Logger())
    router.Use(middleware.CORS())

    // Health check (no auth required)
    router.GET("/health", HealthCheck)

    // API v1 group
    v1 := router.Group("/api/v1")

    // Auth routes (no authentication required)
    authHandler := NewAuthHandler(authService)
    auth := v1.Group("/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }

    // User routes (authentication required)
    userHandler := NewUserHandler(userService)
    users := v1.Group("/users")
    users.Use(middleware.AuthMiddleware(jwtSecret)) // Apply auth middleware
    {
        users.GET("", userHandler.GetUsers)
        users.GET("/:id", userHandler.GetUser)
        users.PUT("/:id", userHandler.UpdateUser)
        users.DELETE("/:id", userHandler.DeleteUser)
    }

    // Crypto routes (authentication required)
    cryptoHandler := NewCryptoHandler(cryptoService)
    crypto := v1.Group("/crypto")
    crypto.Use(middleware.AuthMiddleware(jwtSecret))
    {
        // Standard REST endpoints
        crypto.GET("/:coinId", cryptoHandler.GetSingleCrypto)
        crypto.POST("/bulk", cryptoHandler.GetBulkCrypto)
        crypto.POST("/portfolio", cryptoHandler.GetPortfolioRealtime)
        crypto.GET("/popular", cryptoHandler.GetPopularCoins)
        
        // Cache management
        crypto.GET("/cache/stats", cryptoHandler.GetCacheStats)
        crypto.DELETE("/cache", cryptoHandler.ClearCache)
        
        // Streaming endpoints
        crypto.GET("/stream/prices", cryptoHandler.StreamPrices)
        crypto.GET("/stream/ws", cryptoHandler.WebSocketHandler)
        crypto.POST("/stream/portfolio", cryptoHandler.StreamPortfolio)
    }

    return router
}
```

**Routing Concepts**:
- **Route Groups**: Organize related endpoints with common middleware
- **Middleware Stacking**: Apply middleware at different levels
- **Dependency Injection**: Pass services to handlers through constructors

### Middleware Execution Order
```
Request
   ↓
Logger Middleware
   ↓
CORS Middleware
   ↓
Route Group Middleware (Auth)
   ↓
Handler Function
   ↓
Response
```

## 🔌 API Endpoints Reference

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
}
```

### Cryptocurrency Endpoints

All crypto endpoints require authentication via `Authorization: Bearer <token>` header.

#### Get Single Cryptocurrency
```http
GET /api/v1/crypto/bitcoin
Authorization: Bearer <your-jwt-token>
```

#### Get Popular Cryptocurrencies
```http
GET /api/v1/crypto/popular?limit=5
Authorization: Bearer <your-jwt-token>
```

#### Bulk Cryptocurrency Data
```http
POST /api/v1/crypto/bulk
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "coins": ["bitcoin", "ethereum", "cardano"],
  "timeout": 10
}
```

#### Portfolio Tracking
```http
POST /api/v1/crypto/portfolio
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "coins": ["bitcoin", "ethereum", "solana"]
}
```

### Streaming Endpoints

#### Server-Sent Events
```http
GET /api/v1/crypto/stream/prices?coins=bitcoin,ethereum&interval=5&max_updates=10
Authorization: Bearer <your-jwt-token>
```

#### WebSocket Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/crypto/stream/ws');
// Add Authorization header through subprotocols or query params
```

#### Portfolio Streaming
```http
POST /api/v1/crypto/stream/portfolio
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "coins": ["bitcoin", "ethereum"]
}
```

### Cache Management

#### Get Cache Statistics
```http
GET /api/v1/crypto/cache/stats
Authorization: Bearer <your-jwt-token>
```

#### Clear Cache
```http
DELETE /api/v1/crypto/cache
Authorization: Bearer <your-jwt-token>
```

## 🗄️ Database Integration

### GORM (Go Object-Relational Mapping)

#### Database Connection
```go
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
```

#### Auto Migration
```go
// In main.go
if err := db.AutoMigrate(&models.User{}); err != nil {
    log.Fatal("Failed to migrate database:", err)
}
```

**Auto Migration Benefits**:
- Automatically creates/updates database tables
- Manages schema changes during development
- Handles relationships and indexes

#### GORM Query Examples
```go
// Create
user := models.User{Username: "john", Email: "john@example.com"}
db.Create(&user)

// Read
var user models.User
db.First(&user, 1) // Find by primary key
db.Where("email = ?", "john@example.com").First(&user)

// Update
db.Model(&user).Update("username", "john_updated")
db.Model(&user).Updates(models.User{Username: "john", Email: "new@email.com"})

// Delete
db.Delete(&user, 1)
```

### Database Schema
```sql
-- Users table (auto-generated from User model)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL -- For soft deletes
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

## 🧪 Testing Your Application

### Using curl Commands

#### 1. Register and Login
```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "crypto_trader",
    "email": "trader@example.com",
    "password": "SecurePass123!"
  }'

# Login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "trader@example.com",
    "password": "SecurePass123!"
  }'
```

#### 2. Test Crypto Endpoints
```bash
# Get single crypto data
curl -X GET http://localhost:8080/api/v1/crypto/bitcoin \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get popular cryptos
curl -X GET "http://localhost:8080/api/v1/crypto/popular?limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Bulk crypto data (demonstrates concurrency)
curl -X POST http://localhost:8080/api/v1/crypto/bulk \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "coins": ["bitcoin", "ethereum", "cardano", "solana"],
    "timeout": 15
  }'
```

#### 3. Test Streaming (Server-Sent Events)
```bash
# Stream crypto prices (use -N for no buffering)
curl -N -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  "http://localhost:8080/api/v1/crypto/stream/prices?coins=bitcoin,ethereum&interval=3&max_updates=5"
```

### Using Postman

1. **Set Environment Variables:**
   - `base_url`: `http://localhost:8080`
   - `jwt_token`: (set after login)

2. **Authentication Workflow:**
   - Use register/login endpoints to get JWT token
   - Add to environment variable for reuse

3. **Testing Streaming:**
   - SSE endpoints work directly in Postman
   - See real-time data streaming in response

### WebSocket Testing

#### Browser Console Method
```javascript
// Open browser console and run:
const ws = new WebSocket('ws://localhost:8080/api/v1/crypto/stream/ws');

ws.onopen = function() {
    console.log('Connected to WebSocket');
    ws.send(JSON.stringify({
        action: 'subscribe',
        id: 'test-client-123'
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received crypto update:', data);
};

ws.onclose = function() {
    console.log('WebSocket connection closed');
};

// Send ping to test bidirectional communication
ws.send(JSON.stringify({
    action: 'ping',
    id: 'ping-test'
}));
```

## 🔍 Go Language Learning Points

### Error Handling Patterns
```go
// Standard error handling
result, err := someFunction()
if err != nil {
    log.Printf("Error: %v", err)
    return nil, fmt.Errorf("operation failed: %w", err)
}

// Error wrapping for context
if err != nil {
    return fmt.Errorf("failed to fetch crypto data for %s: %w", coinID, err)
}

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
```

### Interface Usage
```go
// Define behavior, not implementation
type CryptoFetcher interface {
    GetCrypto(coinID string) (*models.CryptoData, error)
    GetBulkCrypto(coins []string) (*models.PortfolioResponse, error)
}

// Multiple implementations
type CoinGeckoFetcher struct{}
type CoinMarketCapFetcher struct{}

// Both implement CryptoFetcher interface
func (c CoinGeckoFetcher) GetCrypto(coinID string) (*models.CryptoData, error) {
    // Implementation specific to CoinGecko
}
```

### Pointer vs Value Semantics
```go
// When to use pointers vs values
type User struct {
    ID   uint
    Name string
}

// Method with pointer receiver (can modify)
func (u *User) UpdateName(name string) {
    u.Name = name
}

// Method with value receiver (read-only)
func (u User) GetDisplayName() string {
    return fmt.Sprintf("User: %s", u.Name)
}

// Function parameters
func ProcessUser(u *User) { /* Can modify original */ }
func DisplayUser(u User)  { /* Works with copy */ }
```

### Package Organization Best Practices
```go
// Good: Organize by feature/domain
internal/
├── auth/          # Authentication domain
│   ├── handler.go
│   ├── service.go
│   └── models.go
├── crypto/        # Cryptocurrency domain
│   ├── handler.go
│   ├── service.go
│   └── models.go

// Avoid: Organizing by technical layer
internal/
├── handlers/      # All handlers together
├── services/      # All services together
├── models/        # All models together
```

## 🚀 Performance Considerations

### Concurrency Best Practices

#### 1. Goroutine Lifecycle Management
```go
// Always ensure goroutines can exit
func (s *Service) StartWorker(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                // Do work
            case <-ctx.Done():
                log.Println("Worker shutting down")
                return // Exit goroutine
            }
        }
    }()
}
```

#### 2. Channel Buffer Sizing
```go
// Unbuffered: Synchronous communication
ch := make(chan int)

// Small buffer: Usually sufficient
ch := make(chan int, 10)

// Buffer size = number of producers: Prevents blocking
ch := make(chan int, numProducers)
```

#### 3. Resource Pooling
```go
// HTTP client reuse
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}
```

### Memory Management

#### 1. Avoid Memory Leaks
```go
// Bad: Goroutine leak
go func() {
    for {
        // This goroutine never exits!
        time.Sleep(time.Second)
    }
}()

// Good: Cancellable goroutine
go func() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Do work
        case <-ctx.Done():
            return // Properly exit
        }
    }
}()
```

#### 2. Efficient String Operations
```go
// Bad: Repeated string concatenation
var result string
for _, item := range items {
    result += item // Creates new string each time
}

// Good: Use strings.Builder
var builder strings.Builder
for _, item := range items {
    builder.WriteString(item)
}
result := builder.String()
```

### Database Performance

#### 1. Connection Pooling
```go
// Configure GORM connection pool
sqlDB, err := db.DB()
if err != nil {
    return err
}

sqlDB.SetMaxIdleConns(10)           // Idle connections
sqlDB.SetMaxOpenConns(100)          // Max connections
sqlDB.SetConnMaxLifetime(time.Hour) // Connection lifetime
```

#### 2. Query Optimization
```go
// Bad: N+1 query problem
var users []User
db.Find(&users)
for _, user := range users {
    var posts []Post
    db.Where("user_id = ?", user.ID).Find(&posts) // N queries
}

// Good: Preload relationships
var users []User
db.Preload("Posts").Find(&users) // 2 queries total
```

## 🛡️ Security Considerations

### JWT Security
```go
// Use strong secrets (in production, use environment variables)
jwtSecret := os.Getenv("JWT_SECRET")
if len(jwtSecret) < 32 {
    log.Fatal("JWT secret must be at least 32 characters")
}

// Set appropriate expiration
claims := jwt.MapClaims{
    "user_id": userID,
    "exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 hour expiry
    "iat":     time.Now().Unix(),                     // Issued at
}
```

### Input Validation
```go
// Validate and sanitize input
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// Additional custom validation
func validatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
        return errors.New("password must contain uppercase letter")
    }
    if !regexp.MustCompile(`[0-9]`).MatchString(password) {
        return errors.New("password must contain number")
    }
    return nil
}
```

### Rate Limiting
```go
// Implement rate limiting for APIs
func rateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(10), 10) // 10 requests per second
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 📈 Monitoring and Observability

### Logging Best Practices
```go
// Structured logging with context
log.Printf("[%s] User %d accessed crypto data for %s - Response: %dms",
    c.Request.Method,
    userID,
    coinID,
    responseTime.Milliseconds(),
)

// Use log levels
log.Println("INFO: Server started on port 8080")
log.Printf("WARN: High response time: %v", responseTime)
log.Printf("ERROR: Failed to fetch crypto data: %v", err)
```

### Metrics Collection
```go
// Track important metrics
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "endpoint", "status"},
    )
)

func metricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            fmt.Sprintf("%d", c.Writer.Status()),
        ).Observe(duration)
    }
}
```

## 🎯 Next Steps and Extensions

### Suggested Improvements

1. **Add Database Migrations**
   - Use golang-migrate for versioned database changes
   - Separate migration files for better version control

2. **Implement Caching Layer**
   - Add Redis for distributed caching
   - Implement cache invalidation strategies

3. **Add Comprehensive Testing**
   - Unit tests for services
   - Integration tests for handlers
   - Load testing for concurrent operations

4. **Enhance Security**
   - Add refresh token mechanism
   - Implement role-based access control (RBAC)
   - Add API key authentication for external clients

5. **Improve Observability**
   - Add distributed tracing with OpenTelemetry
   - Implement health checks
   - Add Prometheus metrics

6. **Production Readiness**
   - Add graceful shutdown
   - Configuration via environment variables
   - Docker containerization
   - Kubernetes deployment manifests

### Learning Resources

- **Go Documentation**: https://golang.org/doc/
- **Effective Go**: https://golang.org/doc/effective_go
- **Go Concurrency Patterns**: https://talks.golang.org/2012/concurrency.slide
- **Gin Framework**: https://gin-gonic.com/docs/
- **GORM Guide**: https://gorm.io/docs/

---

## 📝 Conclusion

This project demonstrates a production-ready Go backend with:

- **Clean Architecture**: Well-organized, maintainable code structure
- **Concurrency Mastery**: Goroutines, channels, wait groups, and mutexes
- **Real-time Features**: SSE and WebSocket implementations
- **Modern Practices**: JWT auth, middleware, proper error handling
- **Performance Optimization**: Caching, rate limiting, concurrent API calls

**Key Go Concepts Learned**:
- Package organization and visibility rules
- Interface-based design for testability
- Error handling and custom error types
- Concurrent programming patterns
- HTTP server implementation with Gin
- Database integration with GORM
- Real-time communication protocols

This foundation prepares you for building scalable, concurrent backend systems in Go. The patterns and practices demonstrated here apply to much larger, production-scale applications.

