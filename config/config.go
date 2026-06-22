package config

import (
	"os"
	"runtime"
	"time"
)

// LoadConfig loads the default configuration
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Attack: AttackConfig{
			Threads:       runtime.NumCPU() * 2,
			Timeout:       10 * time.Second,
			PPS:           100,
			StealthMode:   false,
			ProxyEnabled:  false,
			ResumeEnabled: false,
			PostLoginEnabled: false,
			InsecureTLS:   false,
		},
		Recon: ReconConfig{
			CheckNLA:      true,
			CheckSSL:      true,
			DetectVersion: true,
			MeasureLatency: true,
		},
		Evasion: EvasionConfig{
			MinDelay:     100,
			MaxDelay:     500,
			AdaptiveRate: false,
		},
	}

	return cfg, nil
}

// SaveConfig saves configuration to file
func (c *Config) SaveConfig(filename string) error {
	// TODO: Implement YAML serialization
	return nil
}

// LoadConfigFromFile loads configuration from file
func LoadConfigFromFile(filename string) (*Config, error) {
	// TODO: Implement YAML deserialization
	return LoadConfig()
}

// GetConfigDir returns the config directory
func GetConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + "/.falcon"
}
