package state

import (
	"encoding/json"
	"falcon/config"
	"fmt"
	"os"
	"time"
)

// Checkpoint represents a save point
type Checkpoint struct {
	Timestamp       time.Time
	TargetsProgress map[string]*TargetProgress
	Results         []*config.Result
	Statistics      Statistics
	LastCredentialIndex int
}

// TargetProgress tracks progress for a single target
type TargetProgress struct {
	IP                  string
	Port                int
	Status              string
	LastCredentialTried string
	Attempts            int
	Successes           int
}

// Statistics holds attack statistics
type Statistics struct {
	TotalAttempts    int
	SuccessfulLogins int
	FailedAttempts   int
	SkippedTargets   int
	ElapsedTime      time.Duration
}

// StateManager manages attack state and checkpoints
type StateManager struct {
	StateFile    string
	Checkpoints  []*Checkpoint
	CurrentState *Checkpoint
}

// NewStateManager creates a new state manager
func NewStateManager(stateFile string) *StateManager {
	return &StateManager{
		StateFile:   stateFile,
		Checkpoints: make([]*Checkpoint, 0),
	}
}

// SaveCheckpoint saves the current state
func (sm *StateManager) SaveCheckpoint(targets []*config.Target, results []*config.Result, stats Statistics) error {
	checkpoint := &Checkpoint{
		Timestamp:       time.Now(),
		TargetsProgress: make(map[string]*TargetProgress),
		Results:         results,
		Statistics:      stats,
	}

	for _, target := range targets {
		key := fmt.Sprintf("%s:%d", target.IP, target.Port)
		checkpoint.TargetsProgress[key] = &TargetProgress{
			IP:     target.IP,
			Port:   target.Port,
			Status: target.Status,
		}
	}

	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	err = os.WriteFile(sm.StateFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	sm.CurrentState = checkpoint
	return nil
}

// LoadCheckpoint loads the last saved state
func (sm *StateManager) LoadCheckpoint() (*Checkpoint, error) {
	data, err := os.ReadFile(sm.StateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var checkpoint Checkpoint
	err = json.Unmarshal(data, &checkpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkpoint: %w", err)
	}

	sm.CurrentState = &checkpoint
	return &checkpoint, nil
}

// HasCheckpoint checks if a checkpoint exists
func (sm *StateManager) HasCheckpoint() bool {
	_, err := os.Stat(sm.StateFile)
	return err == nil
}
