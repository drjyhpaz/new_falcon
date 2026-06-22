package config

import (
	"fmt"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load(".env")

	cfg := &Config{
		Attack: AttackConfig{
			Strategy:         getEnvAttackStrategy("ATTACK_STRATEGY", "spray"),
			Threads:          getEnvInt("ATTACK_THREADS", runtime.NumCPU()*2),
			Timeout:          getEnvDuration("ATTACK_TIMEOUT", 10*time.Second),
			RateLimit:        getEnvInt("ATTACK_RATE_LIMIT", 100),
			LockoutThreshold: getEnvInt("LOCKOUT_THRESHOLD", 3),
			LockoutCooldown:  getEnvDuration("LOCKOUT_COOLDOWN", 60*time.Second),
		},
		RDP: RDPConfig{
			Port:         getEnvUint16("RDP_PORT", 3389),
			InsecureTLS:  getEnvBool("INSECURE_TLS", false),
			NLADetection: getEnvBool("NLA_DETECTION", true),
			PreAuthRecon: getEnvBool("PRE_AUTH_RECON", true),
			ReconTimeout: getEnvDuration("RECON_TIMEOUT", 5*time.Second),
		},
		Stealth: StealthConfig{
			Enabled:      getEnvBool("STEALTH_ENABLED", true),
			JitterMin:    getEnvDuration("STEALTH_MIN_JITTER", 500*time.Millisecond),
			JitterMax:    getEnvDuration("STEALTH_MAX_JITTER", 2*time.Second),
			AdaptiveRate: getEnvBool("STEALTH_ADAPTIVE_RATE", true),
			LowAndSlow:   getEnvBool("STEALTH_LOW_AND_SLOW", false),
		},
		Proxy: ProxyConfig{
			Enabled:  getEnvBool("PROXY_ENABLED", false),
			Type:     getEnv("PROXY_TYPE", "socks5"),
			File:     getEnv("PROXY_FILE", "proxies.txt"),
			Rotation: getEnvBool("PROXY_ROTATION", true),
		},
		State: StateConfig{
			ResumeEnabled:      getEnvBool("RESUME_ENABLED", true),
			CheckpointInterval: getEnvInt("CHECKPOINT_INTERVAL", 100),
		},
		Logging: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "logs/falcon.log"),
		},
	}

	return cfg, nil
}

// Helper functions
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvUint16(key string, defaultVal uint16) uint16 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.ParseUint(valStr, 10, 16); err == nil {
		return uint16(val)
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	// Try parsing as milliseconds
	if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
		return time.Duration(val) * time.Millisecond
	}
	// Try parsing as duration string
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAttackStrategy(key string, defaultVal string) AttackStrategy {
	val := getEnv(key, defaultVal)
	switch val {
	case "spray":
		return PasswordSpray
	case "stuff":
		return CredentialStuff
	case "hybrid":
		return Hybrid
	default:
		return PasswordSpray
	}
}
