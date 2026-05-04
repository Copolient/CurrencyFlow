package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     sync.RWMutex
	rate   rate.Limit
	burst  int
	cleanup *time.Ticker
}

func NewIPRateLimiter(rps int, burst int) *IPRateLimiter {
	rl := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		rate:   rate.Limit(rps),
		burst:  burst,
		cleanup: time.NewTicker(5 * time.Minute),
	}

	go func() {
		for range rl.cleanup.C {
			rl.mu.Lock()
			for ip := range rl.ips {
				delete(rl.ips, ip)
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.ips[ip] = limiter
	}

	return limiter
}

func RateLimit(rps int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rps, rps*2)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
