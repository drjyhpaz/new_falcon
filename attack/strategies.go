package attack

import (
	"falcon/config"
	"sync"
)

import (
	"time"
)

// SprayingStrategy implements password spraying (one password against all users)
type SprayingStrategy struct {
	Targets         []*config.Target
	Credentials     []*config.Credential
	CurrentPassword int
	CurrentTarget   int
	Mutex           sync.Mutex
}

// NewSprayingStrategy creates a new spraying strategy
func NewSprayingStrategy(targets []*config.Target, credentials []*config.Credential) *SprayingStrategy {
	return &SprayingStrategy{
		Targets:     targets,
		Credentials: credentials,
	}
}

// GetNextCredential returns the next credential pair
func (ss *SprayingStrategy) GetNextCredential() *config.Credential {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()

	if ss.CurrentPassword >= len(ss.Credentials) {
		return nil
	}

	cred := ss.Credentials[ss.CurrentPassword]
	ss.CurrentPassword++
	return cred
}

// HasMore checks if there are more credentials
func (ss *SprayingStrategy) HasMore() bool {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()
	return ss.CurrentPassword < len(ss.Credentials)
}

// Reset resets the strategy
func (ss *SprayingStrategy) Reset() {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()
	ss.CurrentPassword = 0
	ss.CurrentTarget = 0
}

// StuffingStrategy implements credential stuffing (all passwords against one user)
type StuffingStrategy struct {
	Targets         []*config.Target
	Credentials     []*config.Credential
	CurrentTarget   int
	CurrentPassword int
	Mutex           sync.Mutex
}

// NewStuffingStrategy creates a new stuffing strategy
func NewStuffingStrategy(targets []*config.Target, credentials []*config.Credential) *StuffingStrategy {
	return &StuffingStrategy{
		Targets:     targets,
		Credentials: credentials,
	}
}

// GetNextCredential returns the next credential pair
func (ss *StuffingStrategy) GetNextCredential() *config.Credential {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()

	if ss.CurrentTarget >= len(ss.Targets) || ss.CurrentPassword >= len(ss.Credentials) {
		return nil
	}

	cred := ss.Credentials[ss.CurrentPassword]
	ss.CurrentPassword++

	if ss.CurrentPassword >= len(ss.Credentials) {
		ss.CurrentPassword = 0
		ss.CurrentTarget++
	}

	return cred
}

// HasMore checks if there are more credentials
func (ss *StuffingStrategy) HasMore() bool {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()
	return ss.CurrentTarget < len(ss.Targets) && ss.CurrentPassword < len(ss.Credentials)
}

// Reset resets the strategy
func (ss *StuffingStrategy) Reset() {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()
	ss.CurrentTarget = 0
	ss.CurrentPassword = 0
}
