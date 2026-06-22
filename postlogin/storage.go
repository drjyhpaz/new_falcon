package postlogin

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SessionStorage manages session persistence
type SessionStorage struct {
	Filename string
}

// NewSessionStorage creates a new session storage
func NewSessionStorage(filename string) *SessionStorage {
	return &SessionStorage{
		Filename: filename,
	}
}

// SaveSessions saves sessions to file
func (ss *SessionStorage) SaveSessions(sessions []*Session) error {
	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sessions: %w", err)
	}

	err = os.WriteFile(ss.Filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write sessions file: %w", err)
	}

	return nil
}

// LoadSessions loads sessions from file
func (ss *SessionStorage) LoadSessions() ([]*Session, error) {
	data, err := os.ReadFile(ss.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read sessions file: %w", err)
	}

	var sessions []*Session
	err = json.Unmarshal(data, &sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sessions: %w", err)
	}

	return sessions, nil
}

// KeepAlive sends keep-alive signals to sessions
func (sm *SessionManager) KeepAlive(id string) error {
	session := sm.GetSession(id)
	if session == nil {
		return fmt.Errorf("session not found")
	}

	session.LastActivity = time.Now()
	return nil
}

// DiscoverNetwork performs network discovery from a session
func (s *Session) DiscoverNetwork() error {
	// TODO: Implement network discovery
	// This would scan for other systems on the network
	s.LastActivity = time.Now()
	return nil
}

// GetSystemInfo returns the system information
func (s *Session) GetSystemInfo() map[string]interface{} {
	return s.SystemInfo
}
