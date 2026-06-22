package proxy

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// ProxyType represents the type of proxy
type ProxyType string

const (
	SOCKS5 ProxyType = "socks5"
	HTTP   ProxyType = "http"
	TOR    ProxyType = "tor"
)

// Proxy represents a proxy server
type Proxy struct {
	Type     ProxyType
	Host     string
	Port     int
	Username string
	Password string
}

// ProxyRotator manages proxy rotation
type ProxyRotator struct {
	Proxies      []*Proxy
	CurrentIndex int
	Enabled      bool
}

// LoadProxies loads proxies from file
func LoadProxies(filename string) ([]*Proxy, error) {
	var proxies []*Proxy

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open proxy file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		proxy, err := ParseProxy(line)
		if err != nil {
			continue
		}
		proxies = append(proxies, proxy)
	}

	return proxies, scanner.Err()
}

// ParseProxy parses a proxy string
func ParseProxy(line string) (*Proxy, error) {
	// Format: type://host:port or type://user:pass@host:port
	parts := strings.Split(line, "://")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid proxy format")
	}

	ptype := ProxyType(strings.ToLower(parts[0]))
	addr := parts[1]

	var user, pass, host string
	var port int

	// Check for credentials
	if strings.Contains(addr, "@") {
		credParts := strings.Split(addr, "@")
		userParts := strings.Split(credParts[0], ":")
		user = userParts[0]
		if len(userParts) > 1 {
			pass = userParts[1]
		}
		addr = credParts[1]
	}

	// Parse host:port
	hostParts := strings.Split(addr, ":")
	host = hostParts[0]
	if len(hostParts) > 1 {
		fmt.Sscanf(hostParts[1], "%d", &port)
	}

	if host == "" || port <= 0 {
		return nil, fmt.Errorf("invalid proxy address")
	}

	return &Proxy{
		Type:     ptype,
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
	}, nil
}

// ValidateProxy checks if a proxy is working
func (p *Proxy) Validate() bool {
	addr := fmt.Sprintf("%s:%d", p.Host, p.Port)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// NewProxyRotator creates a new proxy rotator
func NewProxyRotator(proxies []*Proxy) *ProxyRotator {
	return &ProxyRotator{
		Proxies:      proxies,
		CurrentIndex: 0,
		Enabled:      len(proxies) > 0,
	}
}

// GetNextProxy returns the next proxy in rotation
func (pr *ProxyRotator) GetNextProxy() *Proxy {
	if !pr.Enabled || len(pr.Proxies) == 0 {
		return nil
	}

	proxy := pr.Proxies[pr.CurrentIndex]
	pr.CurrentIndex = (pr.CurrentIndex + 1) % len(pr.Proxies)
	return proxy
}

// ValidateAllProxies validates all proxies
func (pr *ProxyRotator) ValidateAllProxies() int {
	validCount := 0
	for _, proxy := range pr.Proxies {
		if proxy.Validate() {
			validCount++
		}
	}
	return validCount
}
