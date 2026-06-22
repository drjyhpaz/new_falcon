package config

import "time"

// AttackConfig holds the attack configuration
type AttackConfig struct {
	Threads              int
	Timeout              time.Duration
	PPS                  int // Packets Per Second
	StealthMode          bool
	ProxyEnabled         bool
	ProxyType            string // SOCKS5, HTTP, TOR
	ProxyFile            string
	ResumeEnabled        bool
	PostLoginEnabled     bool
	DefaultDomain        string
	InsecureTLS          bool
}

// ReconConfig holds recon settings
type ReconConfig struct {
	CheckNLA      bool
	CheckSSL      bool
	DetectVersion bool
	MeasureLatency bool
}

// EvasionConfig holds evasion settings
type EvasionConfig struct {
	MinDelay    int // milliseconds
	MaxDelay    int // milliseconds
	AdaptiveRate bool
}

// ProxyConfig holds proxy settings
type ProxyConfig struct {
	Type     string   // SOCKS5, HTTP, TOR
	List     []string // List of proxies
	Rotation bool
}

// Target represents an attack target
type Target struct {
	IP       string
	Port     int
	Service  string
	Version  string
	NLAEnabled bool
	SSLEnabled bool
	Latency  int // milliseconds
	Status   string // open, closed, filtered
}

// Credential represents a username/password pair
type Credential struct {
	Username string
	Password string
	Domain   string
}

// Result represents an attack result
type Result struct {
	IP        string
	Port      int
	Username  string
	Password  string
	Domain    string
	Success   bool
	Timestamp time.Time
	Error     string
	PostLogin map[string]interface{}
}

// Config is the main configuration structure
type Config struct {
	Attack AttackConfig
	Recon  ReconConfig
	Evasion EvasionConfig
	Proxy  ProxyConfig
}
