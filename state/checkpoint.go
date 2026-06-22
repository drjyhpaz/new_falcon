package state

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/falconjonz/falcon_rdp/config"
)

// CheckpointState holds the complete state for resuming attacks
type CheckpointState struct {
	Timestamp            time.Time                 `json:"timestamp"`
	LastCredentialIndex  int                       `json:"last_credential_index"`
	LastTargetIndex      int                       `json:"last_target_index"`
	Targets              []config.Target           `json:"targets"`
	Credentials          []config.Credential       `json:"credentials"`
	SuccessfulLogins     []config.Result           `json:"successful_logins"`
	AttemptedPairs       map[string]bool           `json:"attempted_pairs"`
	TotalAttempts        int64                     `json:"total_attempts"`
	SuccessfulAttempts   int64                     `json:"successful_attempts"`
	FailedAttempts       int64                     `json:"failed_attempts"`
	LockoutTracker       map[string]int            `json:"lockout_tracker"`
	StartTime            time.Time                 `json:"start_time"`
	AttackStrategy       string                    `json:"attack_strategy"`
}

// StateManager manages checkpoints and resuming
type StateManager struct {
	state              *CheckpointState
	statefile          string
	mu                 sync.RWMutex
	checkpointInterval int
	attemptCounter     int64
	autosaveEnabled    bool
}

// NewStateManager creates a new state manager
func NewStateManager(statefile string, checkpointInterval int) *StateManager {
	return &StateManager{
		state: &CheckpointState{
			Timestamp:        time.Now(),
			AttemptedPairs:   make(map[string]bool),
			LockoutTracker:   make(map[string]int),
		},
		statefile:          statefile,
		checkpointInterval: checkpointInterval,
		autosaveEnabled:    true,
	}
}

// SaveCheckpoint saves current state to file
func (sm *StateManager) SaveCheckpoint() error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sm.state.Timestamp = time.Now()

	data, err := json.MarshalIndent(sm.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	if err := os.WriteFile(sm.statefile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	return nil
}

// LoadCheckpoint loads state from file
func (sm *StateManager) LoadCheckpoint() error {
	data, err := os.ReadFile(sm.statefile)
	if err != nil {
		return fmt.Errorf("failed to read state file: %v", err)
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if err := json.Unmarshal(data, sm.state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}

	return nil
}

// SetTargets sets attack targets
func (sm *StateManager) SetTargets(targets []config.Target) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.Targets = targets
}

// SetCredentials sets attack credentials
func (sm *StateManager) SetCredentials(creds []config.Credential) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.Credentials = creds
}

// RecordAttempt records an attempt
func (sm *StateManager) RecordAttempt(target config.Target, cred config.Credential, success bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	key := fmt.Sprintf("%s:%d_%s:%s", target.IP, target.Port, cred.Username, cred.Password)
	sm.state.AttemptedPairs[key] = true

	sm.state.TotalAttempts++
	if success {
		sm.state.SuccessfulAttempts++
	} else {
		sm.state.FailedAttempts++
	}

	sm.attemptCounter++

	// Auto-save at checkpoint interval
	if sm.autosaveEnabled && sm.attemptCounter%int64(sm.checkpointInterval) == 0 {
		_ = sm.SaveCheckpoint()
	}
}

// RecordSuccess records a successful login
func (sm *StateManager) RecordSuccess(result config.Result) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.SuccessfulLogins = append(sm.state.SuccessfulLogins, result)
}

// UpdateLockoutTracker updates lockout tracker
func (sm *StateManager) UpdateLockoutTracker(key string, count int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.LockoutTracker[key] = count
}

// GetLockoutTracker returns lockout tracker
func (sm *StateManager) GetLockoutTracker() map[string]int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Return a copy
	tracker := make(map[string]int)
	for k, v := range sm.state.LockoutTracker {
		tracker[k] = v
	}
	return tracker
}

// IsAttempted checks if a pair has been attempted
func (sm *StateManager) IsAttempted(target config.Target, cred config.Credential) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	key := fmt.Sprintf("%s:%d_%s:%s", target.IP, target.Port, cred.Username, cred.Password)
	return sm.state.AttemptedPairs[key]
}

// GetCheckpointState returns current checkpoint state
func (sm *StateManager) GetCheckpointState() CheckpointState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return *sm.state
}

// SetLastIndices sets last processed indices
func (sm *StateManager) SetLastIndices(targetIdx, credIdx int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.LastTargetIndex = targetIdx
	sm.state.LastCredentialIndex = credIdx
}

// GetLastIndices returns last processed indices
func (sm *StateManager) GetLastIndices() (int, int) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.state.LastTargetIndex, sm.state.LastCredentialIndex
}

// GetSuccessfulLogins returns all successful logins
func (sm *StateManager) GetSuccessfulLogins() []config.Result {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	results := make([]config.Result, len(sm.state.SuccessfulLogins))
	copy(results, sm.state.SuccessfulLogins)
	return results
}

// CheckpointExists checks if checkpoint file exists
func (sm *StateManager) CheckpointExists() bool {
	_, err := os.Stat(sm.statefile)
	return err == nil
}

// DeleteCheckpoint deletes checkpoint file
func (sm *StateManager) DeleteCheckpoint() error {
	return os.Remove(sm.statefile)
}

// Clear clears all state
func (sm *StateManager) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.state = &CheckpointState{
		Timestamp:      time.Now(),
		AttemptedPairs: make(map[string]bool),
		LockoutTracker: make(map[string]int),
	}
	sm.attemptCounter = 0
}
