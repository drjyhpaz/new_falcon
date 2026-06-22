package recon

import (
	"falcon/config"
	"falcon/logger"
	"fmt"
	"net"
	"time"
)

// ScanTarget performs reconnaissance on a target
func ScanTarget(target *config.Target, cfg *config.Config) error {
	logger.Info("Starting recon on %s:%d", target.IP, target.Port)

	// Check if port is open
	if err := checkPort(target); err != nil {
		target.Status = "closed"
		return fmt.Errorf("port closed: %w", err)
	}
	target.Status = "open"

	// Measure latency
	if cfg.Recon.MeasureLatency {
		latency := measureLatency(target)
		target.Latency = latency
		logger.Info("Latency to %s:%d: %dms", target.IP, target.Port, latency)
	}

	return nil
}

func checkPort(target *config.Target) error {
	addr := fmt.Sprintf("%s:%d", target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func measureLatency(target *config.Target) int {
	start := time.Now()
	addr := fmt.Sprintf("%s:%d", target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return -1
	}
	defer conn.Close()
	return int(time.Since(start).Milliseconds())
}

// ScanTargets performs reconnaissance on multiple targets
func ScanTargets(targets []*config.Target, cfg *config.Config) error {
	for _, target := range targets {
		if err := ScanTarget(target, cfg); err != nil {
			logger.Warning("Recon failed for %s:%d: %v", target.IP, target.Port, err)
		}
	}
	return nil
}
