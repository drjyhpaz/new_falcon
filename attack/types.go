package attack

import (
	"falcon/config"
	"sync"
	"time"
)

// AttackMode represents the attack mode
type AttackMode string

const (
	PasswordSpraying AttackMode = "spraying"
	CredentialStuffing AttackMode = "stuffing"
	Hybrid AttackMode = "hybrid"
)

// Strategy defines the attack strategy
type Strategy interface {
	GetNextCredential() *config.Credential
	HasMore() bool
	Reset()
}

// AttackEngine manages the brute-force attack
type AttackEngine struct {
	Config          *config.Config
	Targets         []*config.Target
	Credentials     []*config.Credential
	WorkerPool      *WorkerPool
	Results         []*config.Result
	ResultsMutex    sync.RWMutex
	Mode            AttackMode
	Strategy        Strategy
	Running         bool
	Stopped         chan bool
	RateLimiter     *RateLimiter
	LockoutManager  *LockoutManager
}

// LockoutManager handles lockout prevention
type LockoutManager struct {
	FailedAttempts map[string]int // key: "ip:port:username"
	Cooldown       map[string]time.Time
	Threshold      int
	CooldownDuration time.Duration
	Mutex          sync.RWMutex
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	Tokens   int
	Capacity int
	Rate     int // tokens per second
	LastFill time.Time
	Mutex    sync.Mutex
}

// WorkerPool manages concurrent workers
type WorkerPool struct {
	NumWorkers int
	Jobs       chan *Job
	Results    chan *config.Result
	Wg         sync.WaitGroup
	Ctx        interface{} // context for cancellation
}

// Job represents a single attack job
type Job struct {
	Target     *config.Target
	Credential *config.Credential
	Mode       AttackMode
}
