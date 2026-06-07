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
	done   chan struct{}
}

func NewIPRateLimiter(rps int, burst int) *IPRateLimiter {
	rl := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		rate:   rate.Limit(rps),
		burst:  burst,
		cleanup: time.NewTicker(5 * time.Minute),
		done:   make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-rl.done:
				rl.cleanup.Stop()
				return
			case <-rl.cleanup.C:
				rl.mu.Lock()
				for ip := range rl.ips {
					delete(rl.ips, ip)
				}
				rl.mu.Unlock()
			}
		}
	}()

	return rl
}

// Stop shuts down the cleanup goroutine.
func (rl *IPRateLimiter) Stop() {
	close(rl.done)
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

// RateLimit returns a gin.HandlerFunc and the underlying IPRateLimiter.
// Call limiter.Stop() during graceful shutdown to avoid goroutine leaks.
func RateLimit(rps int) (gin.HandlerFunc, *IPRateLimiter) {
	limiter := NewIPRateLimiter(rps, rps*2)

	fn := func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
	return fn, limiter
}
