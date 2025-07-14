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
- **Real-time Streaming**: Server-Sent Events (SSE) and WebSockets with Authentication
- **Database Integration**: PostgreSQL with GORM
- **Authentication**: JWT-based auth system with WebSocket support
- **Clean Architecture**: Separation of concerns with layers
- **Error Handling**: Proper Go error handling patterns
- **Configuration Management**: Environment-based configuration
- **Caching**: In-memory caching with thread safety
- **Rate Limiting**: Controlling external API calls

## 🚀 Quick Start Guide

### Prerequisites

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **PostgreSQL**: Running on configured port (default: 5432)
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
The application is pre-configured with a `.env` file. Update it with your database credentials:

```env
# Server Configuration
PORT=8095
HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapp
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=tHiSiSaSeCrEtKeYfOrJwTtOkEnS
JWT_EXPIRES_IN=24h

# Application Environment
APP_ENV=development
```

4. **Run the application**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8095`

## 🔧 Environment Configuration

The application uses environment variables with sensible defaults. All configuration is managed through the `.env` file:

- **PORT**: Server port (default: 8095)
- **HOST**: Server host (default: localhost)
- **DB_***: Database connection parameters
- **JWT_SECRET**: Secret key for JWT tokens (change in production!)
- **JWT_EXPIRES_IN**: Token expiration time (default: 24h)

**Security Note**: Always use strong, unique JWT secrets in production and never commit sensitive credentials to version control.

## 🏗️ Project Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

```
my-go-backend/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Web server entry point
├── configs/                # Configuration management
│   └── config.go          # Environment-based config
├── internal/               # Private application code
│   ├── handlers/          # HTTP request handlers (controllers)
│   │   ├── auth.go        # Authentication endpoints
│   │   ├── crypto.go      # Cryptocurrency endpoints + WebSocket
│   │   ├── health.go      # Health check endpoint
│   │   ├── routes.go      # Route configuration
│   │   └── user.go        # User management endpoints
│   ├── middleware/        # HTTP middleware
│   │   ├── auth.go        # JWT authentication middleware
│   │   ├── cors.go        # CORS configuration
│   │   └── logger.go      # Request logging
│   └── services/          # Business logic layer
│       ├── auth.go        # Authentication service
│       ├── crypto.go      # Cryptocurrency service
│       └── user.go        # User service
├── pkg/                   # Public library code
│   └── models/           # Data structures and models
│       ├── crypto.go     # Cryptocurrency models
│       ├── response.go   # API response models
│       └── user.go       # User models
├── websocket-test/        # WebSocket testing utilities
│   └── index.html        # Browser-based WebSocket test client
├── .env                   # Environment variables
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md             # This comprehensive guide
```

## 🆕 WebSocket Authentication Feature

This project includes an **advanced WebSocket implementation** with JWT authentication support:

### Authentication Methods

1. **Query Parameter Authentication** (Browser-friendly):
```javascript
const ws = new WebSocket('ws://localhost:8095/api/v1/crypto/stream/ws?token=YOUR_JWT_TOKEN');
```

2. **Authorization Header** (for clients that support it):
```javascript
// Headers can be set via subprotocols in some WebSocket clients
```

### WebSocket Test Client

The project includes a ready-to-use HTML test client at `websocket-test/index.html`:

1. **Open the test client**: Open `websocket-test/index.html` in your browser
2. **Login**: Click "Login" to authenticate and get a JWT token
3. **Connect**: Click "Connect WebSocket" to establish authenticated connection
4. **Test**: Use "Send Ping" and "Send Subscribe" buttons to test bidirectional communication

### WebSocket Implementation Details

```go
// WebSocket handler with JWT authentication
func (h *CryptoHandler) WebSocketHandlerWithAuth(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check for token in query parameter or Authorization header
        tokenString := c.Query("token")
        if tokenString == "" {
            authHeader := c.GetHeader("Authorization")
            if authHeader != "" {
                tokenString = strings.TrimPrefix(authHeader, "Bearer ")
            }
        }

        // Validate JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })

        // Upgrade to WebSocket after authentication
        conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
        // ... handle WebSocket communication
    }
}
```

## 🔌 API Endpoints Reference

### Server Information
- **Base URL**: `http://localhost:8095`
- **API Version**: `v1`
- **Health Check**: `GET /health`

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "crypto_trader",
  "email": "trader@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "crypto_trader",
    "email": "trader@example.com"
  }
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "trader@example.com",
  "password": "SecurePass123!"
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
      "username": "crypto_trader",
      "email": "trader@example.com"
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

#### Bulk Cryptocurrency Data (Demonstrates Concurrency)
```http
POST /api/v1/crypto/bulk
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "coins": ["bitcoin", "ethereum", "cardano", "solana"],
  "timeout": 15
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

### Real-time Streaming Endpoints

#### Server-Sent Events (SSE)
```http
GET /api/v1/crypto/stream/prices?coins=bitcoin,ethereum&interval=5&max_updates=10
Authorization: Bearer <your-jwt-token>
```

#### WebSocket Connection (NEW!)
```javascript
// Method 1: Query parameter (browser-friendly)
const token = "your-jwt-token-here";
const ws = new WebSocket(`ws://localhost:8095/api/v1/crypto/stream/ws?token=${token}`);

ws.onopen = function() {
    console.log('Connected to authenticated WebSocket');
    
    // Send subscription message
    ws.send(JSON.stringify({
        action: 'subscribe',
        id: 'client-123'
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received crypto update:', data);
};

// Test bidirectional communication
ws.send(JSON.stringify({
    action: 'ping',
    id: 'ping-test'
}));
```

**WebSocket Message Types:**
- `ping` → `pong`: Health check
- `subscribe` → `subscribed`: Join crypto updates stream

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

## 🧪 Testing Your Application

### Using curl Commands

#### 1. Register and Login
```bash
# Register a new user
curl -X POST http://localhost:8095/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "crypto_trader",
    "email": "trader@example.com",
    "password": "SecurePass123!"
  }'

# Login to get JWT token
curl -X POST http://localhost:8095/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "trader@example.com",
    "password": "SecurePass123!"
  }'
```

#### 2. Test Crypto Endpoints
```bash
# Get single crypto data
curl -X GET http://localhost:8095/api/v1/crypto/bitcoin \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get popular cryptos
curl -X GET "http://localhost:8095/api/v1/crypto/popular?limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Bulk crypto data (demonstrates concurrency)
curl -X POST http://localhost:8095/api/v1/crypto/bulk \
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
  "http://localhost:8095/api/v1/crypto/stream/prices?coins=bitcoin,ethereum&interval=3&max_updates=5"
```

### WebSocket Testing Options

#### Option 1: Built-in HTML Test Client (Recommended)
1. Start the Go server: `go run cmd/server/main.go`
2. Open `websocket-test/index.html` in your browser
3. Click "Login" to authenticate
4. Click "Connect WebSocket" to establish connection
5. Use "Send Ping" and "Send Subscribe" to test

#### Option 2: Browser Console
```javascript
// First, get a token via fetch or curl, then:
const token = "your-jwt-token-here";
const ws = new WebSocket(`ws://localhost:8095/api/v1/crypto/stream/ws?token=${encodeURIComponent(token)}`);

ws.onopen = () => console.log('Connected!');
ws.onmessage = (event) => console.log('Received:', JSON.parse(event.data));
ws.send(JSON.stringify({ action: 'ping', id: 'test' }));
```

## 🧠 Understanding Go Concurrency in This Project

### Goroutines - Lightweight Threads
```go
// Example from crypto service: Concurrent API calls
for _, coin := range coins {
    wg.Add(1)
    go func(coinID string) {
        defer wg.Done()
        
        // Each coin fetched concurrently
        crypto, err := s.GetSingleCrypto(coinID)
        results <- crypto
    }(coin)
}
```

### Channels - Communication
```go
// Buffered channel for results
results := make(chan models.CryptoData, len(coins))

// Send to channel
results <- cryptoData

// Receive from channel
data := <-results
```

### Mutexes - Thread Safety
```go
// RWMutex for cache access
s.mu.RLock()                    // Multiple readers
if cached, exists := s.cache[coinID]; exists {
    s.mu.RUnlock()
    return &cached, nil
}
s.mu.RUnlock()

s.mu.Lock()                     // Exclusive write
s.cache[coinID] = crypto
s.mu.Unlock()
```

### WebSocket Subscriber Management
```go
// Thread-safe subscriber management for real-time updates
func (s *CryptoService) AddSubscriber(id string) <-chan models.StreamEvent {
    s.subMu.Lock()
    defer s.subMu.Unlock()
    
    eventChan := make(chan models.StreamEvent, 100)
    s.subscribers[id] = eventChan
    return eventChan
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

## 🛡️ Security Features

### JWT Authentication
- **Secure token generation** with configurable expiration
- **Header and query parameter support** for WebSocket compatibility
- **Token validation** on every protected endpoint

### Input Validation
```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### CORS Configuration
- **Cross-origin requests** properly handled
- **Preflight requests** supported for complex requests

### Rate Limiting (via Semaphores)
```go
// Limit concurrent API calls to external service
semaphore := make(chan struct{}, 5) // Max 5 concurrent requests
```

## 🎯 Key Learning Outcomes

After studying this project, you'll understand:

1. **Go Project Organization**: Clean separation between `cmd/`, `internal/`, and `pkg/`
2. **HTTP Server Development**: Gin framework, middleware, and routing
3. **Database Integration**: GORM ORM with PostgreSQL
4. **Concurrency Patterns**: Goroutines, channels, wait groups, and mutexes
5. **Real-time Communication**: SSE and WebSocket implementations
6. **Authentication**: JWT tokens with WebSocket support
7. **Error Handling**: Proper Go error patterns and HTTP status codes
8. **Configuration Management**: Environment-based configuration
9. **Caching Strategies**: Thread-safe in-memory caching
10. **Testing**: Browser-based WebSocket testing

## 🚀 Advanced Features Demonstrated

### 1. Concurrent API Calls with Timeout
```go
// Fetch multiple cryptocurrencies concurrently with timeout
ctx, cancel := context.WithTimeout(context.Background(), timeout)
defer cancel()

// Fan-out pattern: Launch goroutines
for _, coin := range coins {
    wg.Add(1)
    go func(coinID string) {
        defer wg.Done()
        // ... concurrent work with timeout handling
    }(coin)
}
```

### 2. Real-time Data Broadcasting
```go
// Background service that streams price updates
go cryptoService.StartPriceStreaming(ctx, popularCoins, 5*time.Second)

// WebSocket clients receive real-time updates
for event := range eventChan {
    conn.WriteJSON(event)
}
```

### 3. Thread-safe Caching
```go
// Read-heavy cache with RWMutex optimization
func (s *CryptoService) GetCacheStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return map[string]interface{}{
        "cached_coins": len(s.cache),
        "cache_hits":   s.cacheHits,
    }
}
```

### 4. WebSocket Authentication
```go
// Flexible authentication: header or query parameter
tokenString := c.Query("token")
if tokenString == "" {
    authHeader := c.GetHeader("Authorization")
    if authHeader != "" {
        tokenString = strings.TrimPrefix(authHeader, "Bearer ")
    }
}
```

## 🔧 Configuration and Deployment

### Environment Variables
All configuration is externalized via environment variables with sensible defaults:

```bash
# Development
export PORT=8095
export DB_HOST=localhost
export JWT_SECRET=your-secret-key

# Production
export APP_ENV=production
export DB_SSL_MODE=require
export JWT_SECRET=very-secure-production-secret
```

### Docker Support (Future Enhancement)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
CMD ["./main"]
```

## 📈 Performance Considerations

### Optimizations Implemented

1. **Connection Pooling**: HTTP client reuse for external API calls
2. **Concurrent Processing**: Multiple goroutines for bulk operations
3. **Caching Layer**: In-memory cache with TTL (Time To Live)
4. **Rate Limiting**: Semaphore pattern to control external API usage
5. **Efficient JSON Processing**: Streaming JSON for large responses
6. **Channel Buffering**: Prevent blocking in WebSocket broadcasting

### Monitoring Points

- **Response times** for API endpoints
- **Cache hit/miss ratios**
- **Concurrent connection counts**
- **External API rate limit usage**
- **Memory usage** for caching and goroutines

## 🎓 Extended Learning Resources

### Go Fundamentals
- [Effective Go](https://golang.org/doc/effective_go): Official Go best practices
- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide): Advanced concurrency
- [Go by Example](https://gobyexample.com/): Practical Go examples

### Frameworks and Libraries
- [Gin Web Framework](https://gin-gonic.com/docs/): HTTP framework documentation
- [GORM Guide](https://gorm.io/docs/): ORM documentation
- [JWT-Go](https://github.com/golang-jwt/jwt): JWT implementation

### Real-time Communication
- [WebSocket Protocol RFC](https://tools.ietf.org/html/rfc6455): WebSocket specification
- [Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events): SSE documentation

---

## 📝 Conclusion

This Go backend project demonstrates production-ready patterns and practices:

**✅ What You've Built:**
- **Scalable Architecture**: Clean, maintainable code organization
- **Real-time Features**: WebSocket and SSE implementations with authentication
- **Concurrent Processing**: Efficient goroutine patterns for high performance
- **Robust Authentication**: JWT tokens with WebSocket support
- **Comprehensive API**: RESTful endpoints with proper error handling
- **Testing Tools**: Browser-based WebSocket test client

**🎯 Skills Acquired:**
- Go project structure and package organization
- HTTP server development with Gin framework
- Database integration with GORM
- Concurrent programming with goroutines and channels
- Real-time communication protocols (WebSocket/SSE)
- JWT authentication implementation
- Thread-safe caching strategies
- Error handling and logging best practices

**🚀 Next Steps:**
1. Add comprehensive unit and integration tests
2. Implement Redis for distributed caching
3. Add Docker containerization
4. Set up CI/CD pipeline
5. Add monitoring and metrics collection
6. Implement rate limiting middleware
7. Add API documentation with Swagger

This foundation prepares you for building enterprise-scale, concurrent backend systems in Go. The patterns demonstrated here scale from small applications to large, distributed microservices architectures.

**Happy coding! 🚀**