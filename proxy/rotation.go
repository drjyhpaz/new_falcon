package proxy

import (
	"fmt"
	"net"
	"time"
)

// ProxyValidator validates proxy availability and performance
type ProxyValidator struct {
	proxyManager *ProxyManager
	timeout      time.Duration
}

// ProxyValidationResult holds validation results
type ProxyValidationResult struct {
	Proxy    string
	Valid    bool
	Latency  time.Duration
	Error    string
	Timestamp time.Time
}

// NewProxyValidator creates a new proxy validator
func NewProxyValidator(pm *ProxyManager, timeout time.Duration) *ProxyValidator {
	return &ProxyValidator{
		proxyManager: pm,
		timeout:      timeout,
	}
}

// ValidateProxy validates a single proxy
func (pv *ProxyValidator) ValidateProxy(proxy string) ProxyValidationResult {
	result := ProxyValidationResult{
		Proxy:     proxy,
		Timestamp: time.Now(),
	}

	start := time.Now()

	// Try to connect to proxy
	conn, err := net.DialTimeout("tcp", proxy, pv.timeout)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
		pv.proxyManager.MarkProxyFailed(proxy)
		return result
	}
	defer conn.Close()

	result.Valid = true
	result.Latency = time.Since(start)
	pv.proxyManager.MarkProxySuccessful(proxy)

	return result
}

// ValidateAll validates all proxies
func (pv *ProxyValidator) ValidateAll() []ProxyValidationResult {
	pm := pv.proxyManager
	pm.mu.Lock()
	proxyList := make([]string, len(pm.proxies))
	copy(proxyList, pm.proxies)
	pm.mu.Unlock()

	var results []ProxyValidationResult

	for _, proxy := range proxyList {
		result := pv.ValidateProxy(proxy)
		results = append(results, result)

		if result.Valid {
			fmt.Printf("✓ Proxy valid: %s (latency: %dms)\n", proxy, result.Latency.Milliseconds())
		} else {
			fmt.Printf("✗ Proxy invalid: %s (%s)\n", proxy, result.Error)
		}
	}

	return results
}

// GetValidProxyCount returns count of valid proxies
func (pv *ProxyValidator) GetValidProxyCount() int {
	count := 0
	pm := pv.proxyManager
	pm.mu.Lock()
	for proxy := range pm.failedProxies {
		if pm.failedProxies[proxy] == 0 {
			count++
		}
	}
	pm.mu.Unlock()
	return count
}
