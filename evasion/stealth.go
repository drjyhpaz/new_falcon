package evasion

import (
	"math/rand"
	"time"
)

// Stealth manages stealth mode operations
type Stealth struct {
	Enabled      bool
	MinDelay     int // milliseconds
	MaxDelay     int // milliseconds
	AdaptiveRate bool
	CurrentRate  float64
}

// NewStealth creates a new stealth instance
func NewStealth(minDelay, maxDelay int, adaptiveRate bool) *Stealth {
	return &Stealth{
		Enabled:      true,
		MinDelay:     minDelay,
		MaxDelay:     maxDelay,
		AdaptiveRate: adaptiveRate,
		CurrentRate:  1.0,
	}
}

// ApplyJitter adds random delay to requests
func (s *Stealth) ApplyJitter() {
	if !s.Enabled {
		return
	}

	delay := rand.Intn(s.MaxDelay-s.MinDelay) + s.MinDelay
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// AdjustRate reduces the rate based on errors
func (s *Stealth) AdjustRate(hasError bool) {
	if !s.AdaptiveRate {
		return
	}

	if hasError {
		s.CurrentRate *= 0.8 // Reduce by 20%
	} else {
		s.CurrentRate = min(s.CurrentRate*1.05, 1.0) // Increase by 5%, max 1.0
	}
}

// GetDelayMultiplier returns the current delay multiplier
func (s *Stealth) GetDelayMultiplier() float64 {
	return 1.0 / s.CurrentRate
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
