package attack

import (
	"fmt"
	"sync"
	"time"
)

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{
		Capacity: rps,
		Tokens:   rps,
		Rate:     rps,
		LastFill: time.Now(),
	}
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	for rl.Tokens <= 0 {
		rl.Mutex.Unlock()
		time.Sleep(10 * time.Millisecond)
		rl.Mutex.Lock()
		rl.refill()
	}

	rl.Tokens--
}

// refill adds tokens based on elapsed time
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.LastFill)
	tokens := int(elapsed.Seconds()) * rl.Rate
	rl.Tokens = min(rl.Tokens+tokens, rl.Capacity)
	rl.LastFill = now
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SetRate changes the rate
func (rl *RateLimiter) SetRate(rps int) {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	rl.Rate = rps
	rl.Capacity = rps
}

// GetCurrentRate returns the current rate
func (rl *RateLimiter) GetCurrentRate() int {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	return rl.Rate
}
