package evasion

import (
	"math/rand"
	"sync"
	"time"
)

// JitterConfig holds jitter configuration
type JitterConfig struct {
	Minimum time.Duration
	Maximum time.Duration
	Type    JitterType
}

// JitterType defines the type of jitter distribution
type JitterType int

const (
	// UniformJitter: random delay between min and max
	UniformJitter JitterType = iota
	// LinearJitter: gradually increases from min to max
	LinearJitter
	// ExponentialJitter: exponential distribution
	ExponentialJitter
)

// JitterManager manages jitter calculations
type JitterManager struct {
	config   JitterConfig
	counter  int64
	mu       sync.RWMutex
	lastTime time.Time
}

// NewJitterManager creates a new jitter manager
func NewJitterManager(config JitterConfig) *JitterManager {
	return &JitterManager{
		config:   config,
		lastTime: time.Now(),
	}
}

// GetJitter returns a jitter delay based on configuration
func (jm *JitterManager) GetJitter() time.Duration {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	switch jm.config.Type {
	case UniformJitter:
		return jm.uniformJitter()
	case LinearJitter:
		return jm.linearJitter()
	case ExponentialJitter:
		return jm.exponentialJitter()
	default:
		return jm.uniformJitter()
	}
}

// uniformJitter returns random duration between min and max
func (jm *JitterManager) uniformJitter() time.Duration {
	min := jm.config.Minimum.Milliseconds()
	max := jm.config.Maximum.Milliseconds()

	if min == max {
		return jm.config.Minimum
	}

	jitter := rand.Int64N(max-min) + min
	return time.Duration(jitter) * time.Millisecond
}

// linearJitter increases from min to max over time
func (jm *JitterManager) linearJitter() time.Duration {
	jm.counter++

	// Map counter to a value between 0 and 1 (wraps every 1000 requests)
	normalized := float64(jm.counter%1000) / 1000.0

	min := float64(jm.config.Minimum.Milliseconds())
	max := float64(jm.config.Maximum.Milliseconds())

	value := min + (max-min)*normalized
	return time.Duration(int64(value)) * time.Millisecond
}

// exponentialJitter uses exponential distribution
func (jm *JitterManager) exponentialJitter() time.Duration {
	// Exponential distribution: -lambda * ln(random)
	lambda := 1.0 / float64(jm.config.Maximum.Milliseconds())
	random := rand.Float64()
	if random == 0 {
		random = 0.01
	}

	value := -lambda * math.Ln(random)
	if value < float64(jm.config.Minimum.Milliseconds()) {
		value = float64(jm.config.Minimum.Milliseconds())
	}
	if value > float64(jm.config.Maximum.Milliseconds()) {
		value = float64(jm.config.Maximum.Milliseconds())
	}

	return time.Duration(int64(value)) * time.Millisecond
}

// Reset resets jitter state
func (jm *JitterManager) Reset() {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	jm.counter = 0
	jm.lastTime = time.Now()
}
