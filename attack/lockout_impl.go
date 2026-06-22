package attack

import (
	"fmt"
	"sync"
	"time"
)

// NewLockoutManager creates a new lockout manager
func NewLockoutManager(threshold int, cooldownDuration time.Duration) *LockoutManager {
	return &LockoutManager{
		FailedAttempts:   make(map[string]int),
		Cooldown:         make(map[string]time.Time),
		Threshold:        threshold,
		CooldownDuration: cooldownDuration,
	}
}

// RecordFailure records a failed attempt
func (lm *LockoutManager) RecordFailure(key string) {
	lm.Mutex.Lock()
	defer lm.Mutex.Unlock()

	lm.FailedAttempts[key]++
	if lm.FailedAttempts[key] >= lm.Threshold {
		lm.Cooldown[key] = time.Now().Add(lm.CooldownDuration)
	}
}

// IsLockedOut checks if an account is locked out
func (lm *LockoutManager) IsLockedOut(key string) bool {
	lm.Mutex.RLock()
	defer lm.Mutex.RUnlock()

	cooldownTime, exists := lm.Cooldown[key]
	if !exists {
		return false
	}

	if time.Now().After(cooldownTime) {
		delete(lm.Cooldown, key)
		return false
	}

	return true
}

// Reset resets the lockout manager
func (lm *LockoutManager) Reset() {
	lm.Mutex.Lock()
	defer lm.Mutex.Unlock()

	lm.FailedAttempts = make(map[string]int)
	lm.Cooldown = make(map[string]time.Time)
}

// GetFailureCount returns the failure count for a key
func (lm *LockoutManager) GetFailureCount(key string) int {
	lm.Mutex.RLock()
	defer lm.Mutex.RUnlock()
	return lm.FailedAttempts[key]
}

// GetRemainingCooldown returns the remaining cooldown time
func (lm *LockoutManager) GetRemainingCooldown(key string) time.Duration {
	lm.Mutex.RLock()
	defer lm.Mutex.RUnlock()

	cooldownTime, exists := lm.Cooldown[key]
	if !exists {
		return 0
	}

	remaining := cooldownTime.Sub(time.Now())
	if remaining < 0 {
		return 0
	}

	return remaining
}
