package attack

import (
	"falcon/config"
	"fmt"
	"sync"
	"time"
)

// Statistics holds attack statistics
type Statistics struct {
	TotalAttempts    int
	SuccessfulLogins int
	FailedAttempts   int
	SkippedTargets   int
	ElapsedTime      time.Duration
	StartTime        time.Time
	EndTime          time.Time
	CurrentPPS       int
}

// StatisticsManager manages attack statistics
type StatisticsManager struct {
	Stats    Statistics
	Mutex    sync.RWMutex
	StartCh  chan bool
	UpdateCh chan int
}

// NewStatisticsManager creates a new statistics manager
func NewStatisticsManager() *StatisticsManager {
	return &StatisticsManager{
		Stats: Statistics{
			StartTime: time.Now(),
		},
		StartCh:  make(chan bool),
		UpdateCh: make(chan int, 100),
	}
}

// RecordAttempt records an attempt
func (sm *StatisticsManager) RecordAttempt() {
	sm.Mutex.Lock()
	defer sm.Mutex.Unlock()
	sm.Stats.TotalAttempts++
}

// RecordSuccess records a successful login
func (sm *StatisticsManager) RecordSuccess() {
	sm.Mutex.Lock()
	defer sm.Mutex.Unlock()
	sm.Stats.SuccessfulLogins++
}

// RecordFailure records a failed attempt
func (sm *StatisticsManager) RecordFailure() {
	sm.Mutex.Lock()
	defer sm.Mutex.Unlock()
	sm.Stats.FailedAttempts++
}

// GetStatistics returns current statistics
func (sm *StatisticsManager) GetStatistics() Statistics {
	sm.Mutex.RLock()
	defer sm.Mutex.RUnlock()

	st := sm.Stats
	st.ElapsedTime = time.Since(st.StartTime)
	return st
}

// GetSuccessRate returns success rate as percentage
func (sm *StatisticsManager) GetSuccessRate() float64 {
	sm.Mutex.RLock()
	defer sm.Mutex.RUnlock()

	if sm.Stats.TotalAttempts == 0 {
		return 0
	}

	return float64(sm.Stats.SuccessfulLogins) / float64(sm.Stats.TotalAttempts) * 100
}

// GetElapsedTime returns elapsed time
func (sm *StatisticsManager) GetElapsedTime() time.Duration {
	sm.Mutex.RLock()
	defer sm.Mutex.RUnlock()
	return time.Since(sm.Stats.StartTime)
}

// FormatStatistics returns formatted statistics string
func (sm *StatisticsManager) FormatStatistics() string {
	st := sm.GetStatistics()
	return fmt.Sprintf(`
╔════════════════════════════════════════════════════════════╗
║                  ATTACK STATISTICS                         ║
╠════════════════════════════════════════════════════════════╣
║ Total Attempts:        %-40d ║
║ Successful Logins:     %-40d ║
║ Failed Attempts:       %-40d ║
║ Skipped Targets:       %-40d ║
║ Success Rate:          %-39.2f%% ║
║ Elapsed Time:          %-40s ║
║ Current PPS:           %-40d ║
╚════════════════════════════════════════════════════════════╝
`, st.TotalAttempts, st.SuccessfulLogins, st.FailedAttempts,
		st.SkippedTargets, sm.GetSuccessRate(), st.ElapsedTime.String(), st.CurrentPPS)
}
