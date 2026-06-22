package rdp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/falconjonz/falcon_rdp/config"
	"github.com/falconjonz/falcon_rdp/logger"
)

// ReconEngine performs pre-attack reconnaissance
type ReconEngine struct {
	detector       *Detector
	cache          *ReconCache
	log            *logger.Logger
	cfg            *config.Config
	results        []ReconInfo
	mu             sync.RWMutex
	scannedCount   atomic.Int32
	onlineCount    atomic.Int32
	offlineCount   atomic.Int32
}

// NewReconEngine creates a new reconnaissance engine
func NewReconEngine(cfg *config.Config, log *logger.Logger) *ReconEngine {
	return &ReconEngine{
		detector: NewDetector(cfg.RDP.ReconTimeout),
		cache:    NewReconCache(1 * time.Hour),
		log:      log,
		cfg:      cfg,
	}
}

// ScanTargets performs reconnaissance on all targets
func (re *ReconEngine) ScanTargets(targets []config.Target) []ReconInfo {
	re.log.Info("Starting pre-attack reconnaissance...")

	var ips []string
	for _, target := range targets {
		ips = append(ips, target.IP)
	}

	// Determine workers (use half of attack threads)
	workers := re.cfg.Attack.Threads / 2
	if workers < 1 {
		workers = 1
	}

	start := time.Now()
	results := re.detector.ScanMultiple(ips, re.cfg.RDP.Port, workers)

	re.mu.Lock()
	re.results = results
	re.mu.Unlock()

	// Update statistics
	for _, result := range results {
		re.scannedCount.Add(1)
		if result.Open {
			re.onlineCount.Add(1)
		} else {
			re.offlineCount.Add(1)
		}

		// Cache result
		key := fmt.Sprintf("%s:%d", result.IP, result.Port)
		re.cache.Set(key, result)

		// Log result
		if result.Open {
			re.log.Infof("✓ %s:%d (latency: %dms)", result.IP, result.Port, result.Latency.Milliseconds())
			if result.NLAEnabled {
				re.log.Debugf("  └─ NLA enabled")
			}
			if result.SSLEnabled {
				re.log.Debugf("  └─ SSL enabled")
			}
		} else {
			re.log.Debugf("✗ %s:%d (%s)", result.IP, result.Port, result.Error)
		}
	}

	duration := time.Since(start)
	re.log.Infof("Reconnaissance complete: %d targets scanned in %v", len(results), duration)
	re.log.Infof("Online: %d, Offline: %d", re.onlineCount.Load(), re.offlineCount.Load())

	return results
}

// GetOnlineTargets returns only online targets
func (re *ReconEngine) GetOnlineTargets() []ReconInfo {
	re.mu.RLock()
	defer re.mu.RUnlock()

	var onlineTargets []ReconInfo
	for _, result := range re.results {
		if result.Open {
			onlineTargets = append(onlineTargets, result)
		}
	}
	return onlineTargets
}

// GetResults returns all reconnaissance results
func (re *ReconEngine) GetResults() []ReconInfo {
	re.mu.RLock()
	defer re.mu.RUnlock()

	results := make([]ReconInfo, len(re.results))
	copy(results, re.results)
	return results
}

// GetStats returns reconnaissance statistics
func (re *ReconEngine) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"scanned": re.scannedCount.Load(),
		"online":  re.onlineCount.Load(),
		"offline": re.offlineCount.Load(),
	}
}

// CacheResults saves reconnaissance results to cache file
func (re *ReconEngine) CacheResults(filename string) error {
	return re.cache.SaveToFile(filename)
}
