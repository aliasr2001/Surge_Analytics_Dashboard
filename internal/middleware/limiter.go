package middleware

import (
	// "net/http"
	"sync"
	"golang.org/x/time/rate"
)

// IPRateLimiter holds a map of limiters for each IP address
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex // This is a 'Lock' to prevent different threads from crashing the map
}

// NewIPRateLimiter creates the bodyguard
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
	}
}

// GetLimiter returns a token bucket for a specific IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		// Create a new bucket: 'r' is how fast it refills, 'b' is the max size
		limiter = rate.NewLimiter(rate.Limit(1), 5) 
		i.ips[ip] = limiter
	}

	return limiter
}