package middleware

import (
    "net/http"
    "sync"
    "time"
    "go-gin-backend/pkg/response"

    "github.com/gin-gonic/gin"
)

type RateLimiter struct {
    clients map[string]*Client
    mutex   sync.RWMutex
    rate    time.Duration
    burst   int
}

type Client struct {
    tokens    int
    lastToken time.Time
}

func NewRateLimiter(rps int) *RateLimiter {
    return &RateLimiter{
        clients: make(map[string]*Client),
        rate:    time.Second / time.Duration(rps),
        burst:   rps,
    }
}

func (rl *RateLimiter) Allow(clientID string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()

    client, exists := rl.clients[clientID]
    if !exists {
        client = &Client{
            tokens:    rl.burst,
            lastToken: time.Now(),
        }
        rl.clients[clientID] = client
    }

    now := time.Now()
    elapsed := now.Sub(client.lastToken)
    tokensToAdd := int(elapsed / rl.rate)

    if tokensToAdd > 0 {
        client.tokens += tokensToAdd
        if client.tokens > rl.burst {
            client.tokens = rl.burst
        }
        client.lastToken = now
    }

    if client.tokens > 0 {
        client.tokens--
        return true
    }

    return false
}

func RateLimit() gin.HandlerFunc {
    limiter := NewRateLimiter(100) // 100 requests per second

    return func(c *gin.Context) {
        clientIP := c.ClientIP()

        if !limiter.Allow(clientIP) {
            response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded")
            c.Abort()
            return
        }

        c.Next()
    }
}
