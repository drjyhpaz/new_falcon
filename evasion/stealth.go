package evasion

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// StealthManager handles all stealth-related operations
type StealthManager struct {
	minJitter     time.Duration
	maxJitter     time.Duration
	adaptiveRate  bool
	lowAndSlow    bool
	currentRate   float64 // requests per second
	rateMu        sync.RWMutex
	errorsInBatch int
	errorMu       sync.RWMutex
}

// NewStealthManager creates a new stealth manager
func NewStealthManager(minJitter, maxJitter time.Duration, adaptive, lowAndSlow bool) *StealthManager {
	return &StealthManager{
		minJitter:    minJitter,
		maxJitter:    maxJitter,
		adaptiveRate: adaptive,
		lowAndSlow:   lowAndSlow,
		currentRate:  100.0, // default 100 RPS
	}
}

// ApplyJitter applies random delay before request
func (sm *StealthManager) ApplyJitter() time.Duration {
	if sm.minJitter == 0 && sm.maxJitter == 0 {
		return 0
	}

	// Calculate random jitter
	min := sm.minJitter.Milliseconds()
	max := sm.maxJitter.Milliseconds()

	if min == max {
		return time.Duration(min) * time.Millisecond
	}

	jitter := rand.Int64N(max-min) + min
	return time.Duration(jitter) * time.Millisecond
}

// GetDelayBetweenRequests calculates delay between requests based on rate limit
func (sm *StealthManager) GetDelayBetweenRequests() time.Duration {
	sm.rateMu.RLock()
	rrate := sm.currentRate
	sm.rateMu.RUnlock()

	if rrate <= 0 {
		rrate = 1
	}

	// delay = 1 / (requests_per_second)
	delay := time.Duration(float64(time.Second) / rrate)

	// Add jitter
	jitter := sm.ApplyJitter()

	return delay + jitter
}

// AdaptiveRateLimit adjusts rate based on errors
func (sm *StealthManager) AdaptiveRateLimit(error bool) {
	if !sm.adaptiveRate {
		return
	}

	sm.errorMu.Lock()
	if error {
		sm.errorsInBatch++
	} else {
		sm.errorsInBatch = 0
	}
	errorsInBatch := sm.errorsInBatch
	sm.errorMu.Unlock()

	sm.rateMu.Lock()
	defer sm.rateMu.Unlock()

	// If too many errors, reduce rate
	if errorsInBatch > 5 {
		// Reduce rate by 20%
		sm.currentRate *= 0.8
		if sm.currentRate < 1 {
			sm.currentRate = 1
		}
	} else if errorsInBatch == 0 {
		// Gradually increase rate back
		sm.currentRate *= 1.05
		if sm.currentRate > 1000 {
			sm.currentRate = 1000
		}
	}
}

// SetRateLimit sets maximum requests per second
func (sm *StealthManager) SetRateLimit(rps float64) {
	sm.rateMu.Lock()
	defer sm.rateMu.Unlock()

	if rps <= 0 {
		rps = 1
	}
	sm.currentRate = rps
}

// GetRateLimit returns current rate limit
func (sm *StealthManager) GetRateLimit() float64 {
	sm.rateMu.RLock()
	defer sm.rateMu.RUnlock()
	return sm.currentRate
}

// LowAndSlowProfile sets a very conservative rate for slow attacks
func (sm *StealthManager) LowAndSlowProfile() {
	sm.rateMu.Lock()
	defer sm.rateMu.Unlock()

	// 5-10 requests per minute
	sm.currentRate = 0.15 // 9 per minute
	sm.minJitter = 3 * time.Second
	sm.maxJitter = 10 * time.Second
}

// NormalProfile sets normal rate for balanced attacks
func (sm *StealthManager) NormalProfile() {
	sm.rateMu.Lock()
	defer sm.rateMu.Unlock()

	sm.currentRate = 100.0
	sm.minJitter = 100 * time.Millisecond
	sm.maxJitter = 500 * time.Millisecond
}

// AggressiveProfile sets high rate for fast attacks
func (sm *StealthManager) AggressiveProfile() {
	sm.rateMu.Lock()
	defer sm.rateMu.Unlock()

	sm.currentRate = 1000.0
	sm.minJitter = 10 * time.Millisecond
	sm.maxJitter = 50 * time.Millisecond
}

// CalculateBackoffDelay calculates exponential backoff delay
func CalculateBackoffDelay(attempt int, baseDelay time.Duration) time.Duration {
	// Exponential backoff: base * (2 ^ attempt) with jitter
	delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay

	// Add random jitter (±50%)
	jitter := time.Duration(rand.Float64()*float64(delay) - float64(delay)/2)

	totalDelay := delay + jitter
	if totalDelay < 0 {
		totalDelay = baseDelay
	}

	return totalDelay
}

// ShouldThrottle determines if request should be throttled
func (sm *StealthManager) ShouldThrottle() bool {
	if sm.lowAndSlow {
		return rand.Float64() < 0.3 // 30% chance to throttle
	}
	return false
}
