package config

import (
	"time"
)

// AttackStrategy defines the attack mode
type AttackStrategy string

const (
	PasswordSpray    AttackStrategy = "spray"
	CredentialStuff  AttackStrategy = "stuff"
	Hybrid           AttackStrategy = "hybrid"
)

// Config holds all configuration for the tool
type Config struct {
	// Attack Settings
	Attack AttackConfig `json:"attack"`

	// RDP Settings
	RDP RDPConfig `json:"rdp"`

	// Stealth Settings
	Stealth StealthConfig `json:"stealth"`

	// Proxy Settings
	Proxy ProxyConfig `json:"proxy"`

	// State Management
	State StateConfig `json:"state"`

	// Logging
	Logging LogConfig `json:"logging"`
}

type AttackConfig struct {
	Strategy        AttackStrategy `json:"strategy"`
	Threads         int            `json:"threads"`
	Timeout         time.Duration  `json:"timeout"`
	RateLimit       int            `json:"rate_limit"`        // requests per second
	LockoutThreshold int           `json:"lockout_threshold"` // failed attempts before lockout
	LockoutCooldown time.Duration  `json:"lockout_cooldown"`  // cooldown period
}

type RDPConfig struct {
	Port              uint16        `json:"port"`
	InsecureTLS       bool          `json:"insecure_tls"`
	NLADetection      bool          `json:"nla_detection"`
	PreAuthRecon      bool          `json:"pre_auth_recon"`
	ReconTimeout      time.Duration `json:"recon_timeout"`
}

type StealthConfig struct {
	Enabled         bool          `json:"enabled"`
	JitterMin       time.Duration `json:"jitter_min"`
	JitterMax       time.Duration `json:"jitter_max"`
	AdaptiveRate    bool          `json:"adaptive_rate"`
	LowAndSlow      bool          `json:"low_and_slow"`
}

type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`     // socks5, http, tor
	File     string `json:"file"`     // proxy list file
	Rotation bool   `json:"rotation"`
}

type StateConfig struct {
	ResumeEnabled      bool `json:"resume_enabled"`
	CheckpointInterval int  `json:"checkpoint_interval"`
}

type LogConfig struct {
	Level string `json:"level"` // debug, info, warn, error
	File  string `json:"file"`
}

// Target represents an RDP target
type Target struct {
	IP       string        `json:"ip"`
	Port     uint16        `json:"port"`
	Hostname string        `json:"hostname,omitempty"`
	OS       string        `json:"os,omitempty"`
	NLAEnabled bool        `json:"nla_enabled,omitempty"`
	SSLEnabled bool        `json:"ssl_enabled,omitempty"`
	Latency  time.Duration `json:"latency,omitempty"`
	Status   string        `json:"status,omitempty"` // online, offline, recon_failed
}

// Credential represents a username:password pair
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain,omitempty"`
}

// Result represents a successful login
type Result struct {
	IP        string            `json:"ip"`
	Port      uint16            `json:"port"`
	Username  string            `json:"username"`
	Password  string            `json:"password"`
	Domain    string            `json:"domain,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	OSInfo    map[string]string `json:"os_info,omitempty"`
	IsAdmin   bool              `json:"is_admin,omitempty"`
}

// Statistics tracks attack progress
type Statistics struct {
	TotalAttempts   int64
	SuccessfulLogins int64
	FailedAttempts  int64
	SkippedTargets  int64
	StartTime       time.Time
	EndTime         time.Time
}
