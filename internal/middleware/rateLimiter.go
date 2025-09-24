package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	Clients map[string]*rate.Limiter
	mu      sync.RWMutex
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		Clients: make(map[string]*rate.Limiter),
		mu:      sync.RWMutex{},
	}
}

func (rt *RateLimiter) TokenBucketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rt.mu.Lock()
		bucket, exists := rt.Clients[ip]
		if !exists {
			bucket = rate.NewLimiter(1, 5)
			rt.Clients[ip] = bucket
		}
		rt.mu.Unlock()

		if !bucket.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
