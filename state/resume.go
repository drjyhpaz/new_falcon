package state

import (
	"fmt"
	"os"
	"sync"

	"github.com/falconjonz/falcon_rdp/config"
	"github.com/falconjonz/falcon_rdp/logger"
)

// ResumeManager handles resuming interrupted attacks
type ResumeManager struct {
	stateMgr  *StateManager
	log       *logger.Logger
	enabled   bool
	mu        sync.RWMutex
}

// ResumeInfo holds information about resumable state
type ResumeInfo struct {
	CanResume          bool
	LastTimestamp      string
	SuccessfulCount    int64
	TotalAttempts      int64
	TargetCount        int
	CredentialCount    int
	EstimatedProgress  float64
}

// NewResumeManager creates a new resume manager
func NewResumeManager(stateMgr *StateManager, log *logger.Logger, enabled bool) *ResumeManager {
	return &ResumeManager{
		stateMgr: stateMgr,
		log:      log,
		enabled:  enabled,
	}
}

// CanResume checks if attack can be resumed
func (rm *ResumeManager) CanResume() bool {
	if !rm.enabled {
		return false
	}
	return rm.stateMgr.CheckpointExists()
}

// GetResumeInfo returns information about resumable state
func (rm *ResumeManager) GetResumeInfo() ResumeInfo {
	info := ResumeInfo{
		CanResume: rm.CanResume(),
	}

	if !info.CanResume {
		return info
	}

	state := rm.stateMgr.GetCheckpointState()
	info.LastTimestamp = state.Timestamp.Format("2006-01-02 15:04:05")
	info.SuccessfulCount = state.SuccessfulAttempts
	info.TotalAttempts = state.TotalAttempts
	info.TargetCount = len(state.Targets)
	info.CredentialCount = len(state.Credentials)

	// Calculate estimated progress
	if info.TargetCount > 0 && info.CredentialCount > 0 {
		totalPairs := int64(info.TargetCount * info.CredentialCount)
		if totalPairs > 0 {
			info.EstimatedProgress = float64(info.TotalAttempts) / float64(totalPairs) * 100
		}
	}

	return info
}

// Resume resumes attack from checkpoint
func (rm *ResumeManager) Resume() ([]config.Target, []config.Credential, error) {
	if !rm.CanResume() {
		return nil, nil, fmt.Errorf("cannot resume: no checkpoint available")
	}

	if err := rm.stateMgr.LoadCheckpoint(); err != nil {
		return nil, nil, fmt.Errorf("failed to load checkpoint: %v", err)
	}

	state := rm.stateMgr.GetCheckpointState()

	rm.log.Infof("Resuming attack from checkpoint (timestamp: %s)", state.Timestamp.Format("2006-01-02 15:04:05"))
	rm.log.Infof("Previous progress: %d/%d attempts, %d successful",
		state.TotalAttempts,
		int64(len(state.Targets))*int64(len(state.Credentials)),
		state.SuccessfulAttempts,
	)

	return state.Targets, state.Credentials, nil
}

// ClearCheckpoint removes checkpoint file
func (rm *ResumeManager) ClearCheckpoint() error {
	if !rm.stateMgr.CheckpointExists() {
		return nil
	}

	if err := rm.stateMgr.DeleteCheckpoint(); err != nil {
		rm.log.Errorf("failed to delete checkpoint: %v", err)
		return err
	}

	rm.log.Info("Checkpoint cleared")
	return nil
}

// Enable enables resume functionality
func (rm *ResumeManager) Enable() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.enabled = true
}

// Disable disables resume functionality
func (rm *ResumeManager) Disable() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.enabled = false
}

// IsEnabled checks if resume is enabled
func (rm *ResumeManager) IsEnabled() bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.enabled
}

// PromptUserForResume prompts user to resume or start fresh
func (rm *ResumeManager) PromptUserForResume() (bool, error) {
	if !rm.CanResume() {
		return false, nil
	}

	info := rm.GetResumeInfo()
	fmt.Println("\n" + "="*60)
	fmt.Println("CHECKPOINT FOUND - Resume previous attack?")
	fmt.Println("="*60)
	fmt.Printf("Last checkpoint: %s\n", info.LastTimestamp)
	fmt.Printf("Previous attempts: %d\n", info.TotalAttempts)
	fmt.Printf("Successful logins: %d\n", info.SuccessfulCount)
	fmt.Printf("Progress: %.2f%%\n", info.EstimatedProgress)
	fmt.Println("="*60)
	fmt.Print("Resume attack? (y/n): ")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}

	return response == "y" || response == "Y", nil
}
