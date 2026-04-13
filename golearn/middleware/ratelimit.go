package middleware

import (
	"net/http"
	"sync"

	"golearn/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter holds a rate.Limiter per IP address.
type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// get returns (creating if needed) the limiter for the given IP.
func (il *ipLimiter) get(ip string) *rate.Limiter {
	il.mu.Lock()
	defer il.mu.Unlock()
	if lim, ok := il.limiters[ip]; ok {
		return lim
	}
	lim := rate.NewLimiter(il.r, il.b)
	il.limiters[ip] = lim
	return lim
}

// RateLimit enforces per-IP request rate limits globally.
func RateLimit(cfg *config.Config) gin.HandlerFunc {
	il := newIPLimiter(rate.Limit(cfg.RateLimit), cfg.BurstLimit)
	return func(c *gin.Context) {
		if !il.get(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests — slow down",
			})
			return
		}
		c.Next()
	}
}
