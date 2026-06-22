package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// ProxyManager manages proxy list and rotation
type ProxyManager struct {
	proxies       []string
	currentIndex  atomic.Int32
	mu            sync.RWMutex
	proxyType     string // socks5, http, tor
	enabled       bool
	rotationMode  RotationMode
	failedProxies map[string]int // track failed proxies
}

// RotationMode defines how proxies are rotated
type RotationMode int

const (
	// RoundRobin: cycle through proxies in order
	RoundRobin RotationMode = iota
	// Random: select random proxy
	Random
	// LeastUsed: select least used proxy
	LeastUsed
)

// NewProxyManager creates a new proxy manager
func NewProxyManager(proxyType string) *ProxyManager {
	return &ProxyManager{
		proxyType:     proxyType,
		rotationMode:  RoundRobin,
		failedProxies: make(map[string]int),
	}
}

// LoadProxies loads proxies from file
func (pm *ProxyManager) LoadProxies(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open proxy file: %v", err)
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			if err := pm.validateProxy(line); err == nil {
				proxies = append(proxies, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading proxy file: %v", err)
	}

	if len(proxies) == 0 {
		return fmt.Errorf("no valid proxies found")
	}

	pm.mu.Lock()
	pm.proxies = proxies
	pm.mu.Unlock()

	return nil
}

// AddProxy adds a single proxy
func (pm *ProxyManager) AddProxy(proxy string) error {
	if err := pm.validateProxy(proxy); err != nil {
		return err
	}

	pm.mu.Lock()
	pm.proxies = append(pm.proxies, proxy)
	pm.mu.Unlock()

	return nil
}

// GetProxy returns next proxy based on rotation mode
func (pm *ProxyManager) GetProxy() (string, error) {
	pm.mu.RLock()
	if len(pm.proxies) == 0 {
		pm.mu.RUnlock()
		return "", fmt.Errorf("no proxies available")
	}
	proxyList := pm.proxies
	pm.mu.RUnlock()

	switch pm.rotationMode {
	case RoundRobin:
		return pm.getRoundRobin(proxyList), nil
	case Random:
		return pm.getRandom(proxyList), nil
	case LeastUsed:
		return pm.getLeastUsed(proxyList), nil
	default:
		return pm.getRoundRobin(proxyList), nil
	}
}

// getRoundRobin returns proxy in round-robin fashion
func (pm *ProxyManager) getRoundRobin(proxies []string) string {
	idx := pm.currentIndex.Add(1) % int32(len(proxies))
	return proxies[idx]
}

// getRandom returns random proxy
func (pm *ProxyManager) getRandom(proxies []string) string {
	// Using a simple pseudo-random selection
	idx := (pm.currentIndex.Add(1) * 7) % int32(len(proxies))
	if idx < 0 {
		idx = -idx
	}
	return proxies[idx]
}

// getLeastUsed returns least used proxy
func (pm *ProxyManager) getLeastUsed(proxies []string) string {
	pm.mu.RLock()
	failedProxies := pm.failedProxies
	pm.mu.RUnlock()

	var leastUsedProxy string
	minFails := int(^uint(0) >> 1) // max int

	for _, proxy := range proxies {
		fails := failedProxies[proxy]
		if fails < minFails {
			minFails = fails
			leastUsedProxy = proxy
		}
	}

	if leastUsedProxy == "" && len(proxies) > 0 {
		return proxies[0]
	}
	return leastUsedProxy
}

// MarkProxyFailed marks a proxy as failed
func (pm *ProxyManager) MarkProxyFailed(proxy string) {
	pm.mu.Lock()
	pm.failedProxies[proxy]++
	pm.mu.Unlock()
}

// MarkProxySuccessful marks a proxy as successful
func (pm *ProxyManager) MarkProxySuccessful(proxy string) {
	pm.mu.Lock()
	if pm.failedProxies[proxy] > 0 {
		pm.failedProxies[proxy]--
	}
	pm.mu.Unlock()
}

// GetProxyURL returns proxy URL for the given proxy string
func (pm *ProxyManager) GetProxyURL(proxy string) (*url.URL, error) {
	var proxyURL string

	switch pm.proxyType {
	case "socks5":
		proxyURL = fmt.Sprintf("socks5://%s", proxy)
	case "http":
		proxyURL = fmt.Sprintf("http://%s", proxy)
	case "tor":
		proxyURL = fmt.Sprintf("socks5://%s", proxy) // TOR uses SOCKS5
	default:
		return nil, fmt.Errorf("unsupported proxy type: %s", pm.proxyType)
	}

	return url.Parse(proxyURL)
}

// validateProxy validates proxy format
func (pm *ProxyManager) validateProxy(proxy string) error {
	parts := strings.Split(proxy, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid proxy format: %s", proxy)
	}

	ip := parts[0]
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP in proxy: %s", ip)
	}

	return nil
}

// Count returns total proxy count
func (pm *ProxyManager) Count() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.proxies)
}

// SetRotationMode sets proxy rotation mode
func (pm *ProxyManager) SetRotationMode(mode RotationMode) {
	pm.rotationMode = mode
}

// Enable enables proxy usage
func (pm *ProxyManager) Enable() {
	pm.enabled = true
}

// Disable disables proxy usage
func (pm *ProxyManager) Disable() {
	pm.enabled = false
}

// IsEnabled returns if proxy is enabled
func (pm *ProxyManager) IsEnabled() bool {
	return pm.enabled
}
