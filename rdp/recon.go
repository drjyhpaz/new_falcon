package rdp

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// ReconCache caches reconnaissance results
type ReconCache struct {
	cache map[string]ReconInfo
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewReconCache creates a new recon cache
func NewReconCache(ttl time.Duration) *ReconCache {
	return &ReconCache{
		cache: make(map[string]ReconInfo),
		ttl:   ttl,
	}
}

// Get retrieves cached recon info
func (rc *ReconCache) Get(key string) (ReconInfo, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	info, exists := rc.cache[key]
	if !exists {
		return ReconInfo{}, false
	}

	// Check if expired
	if time.Since(info.Timestamp) > rc.ttl {
		return ReconInfo{}, false
	}

	return info, true
}

// Set stores recon info in cache
func (rc *ReconCache) Set(key string, info ReconInfo) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.cache[key] = info
}

// Clear clears the cache
func (rc *ReconCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.cache = make(map[string]ReconInfo)
}

// SaveToFile saves cache to JSON file
func (rc *ReconCache) SaveToFile(filename string) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	var results []ReconInfo
	for _, info := range rc.cache {
		results = append(results, info)
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal recon cache: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write recon cache: %v", err)
	}

	return nil
}

// LoadFromFile loads cache from JSON file
func (rc *ReconCache) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read recon cache file: %v", err)
	}

	var results []ReconInfo
	if err := json.Unmarshal(data, &results); err != nil {
		return fmt.Errorf("failed to unmarshal recon cache: %v", err)
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	for _, info := range results {
		key := fmt.Sprintf("%s:%d", info.IP, info.Port)
		rc.cache[key] = info
	}

	return nil
}

// GetStats returns cache statistics
func (rc *ReconCache) GetStats() map[string]interface{} {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	onlineCount := 0
	offlineCount := 0
	nlaCount := 0
	sslCount := 0

	for _, info := range rc.cache {
		if info.Open {
			onlineCount++
			if info.NLAEnabled {
				nlaCount++
			}
			if info.SSLEnabled {
				sslCount++
			}
		} else {
			offlineCount++
		}
	}

	return map[string]interface{}{
		"total":     len(rc.cache),
		"online":    onlineCount,
		"offline":   offlineCount,
		"nla_count": nlaCount,
		"ssl_count": sslCount,
	}
}
