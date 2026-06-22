package rdp

import (
	"fmt"
	"net"
	"time"
)

// ReconInfo holds reconnaissance information about an RDP target
type ReconInfo struct {
	IP           string
	Port         uint16
	Open         bool
	Latency      time.Duration
	NLAEnabled   bool
	SSLEnabled   bool
	WindowsVersion string
	Error        string
	Timestamp    time.Time
}

// Detector performs pre-attack reconnaissance on RDP targets
type Detector struct {
	timeout time.Duration
}

// NewDetector creates a new RDP detector
func NewDetector(timeout time.Duration) *Detector {
	return &Detector{
		timeout: timeout,
	}
}

// DetectRDP checks if RDP is open on target
func (d *Detector) DetectRDP(ip string, port uint16) (bool, time.Duration, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)

	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, d.timeout)
	latency := time.Since(start)

	if err != nil {
		return false, latency, err
	}
	defer conn.Close()

	return true, latency, nil
}

// DetectNLA checks if NLA is enabled
func (d *Detector) DetectNLA(ip string, port uint16) (bool, error) {
	// Connect to RDP port
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, d.timeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// Send X.224 Connection Request
	req := buildX224ConnectionRequest()
	_, err = conn.Write(req)
	if err != nil {
		return false, err
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return false, err
	}

	// Parse response to detect NLA
	nlaEnabled := detectNLAFromResponse(buffer[:n])
	return nlaEnabled, nil
}

// DetectSSL checks if SSL/TLS is enabled
func (d *Detector) DetectSSL(ip string, port uint16) (bool, error) {
	// Connect to RDP port
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, d.timeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// Try SSL/TLS handshake
	// For now, return false as placeholder
	// In production, would use crypto/tls package
	return false, nil
}

// Scan performs comprehensive reconnaissance on target
func (d *Detector) Scan(ip string, port uint16) ReconInfo {
	recon := ReconInfo{
		IP:        ip,
		Port:      port,
		Timestamp: time.Now(),
	}

	// Check if RDP is open
	open, latency, err := d.DetectRDP(ip, port)
	recon.Open = open
	recon.Latency = latency

	if !open {
		recon.Error = fmt.Sprintf("RDP not open: %v", err)
		return recon
	}

	// Detect NLA
	nla, err := d.DetectNLA(ip, port)
	if err == nil {
		recon.NLAEnabled = nla
	}

	// Detect SSL
	ssl, err := d.DetectSSL(ip, port)
	if err == nil {
		recon.SSLEnabled = ssl
	}

	return recon
}

// ScanMultiple performs reconnaissance on multiple targets concurrently
func (d *Detector) ScanMultiple(targets []string, port uint16, workers int) []ReconInfo {
	resultsChan := make(chan ReconInfo, len(targets))
	workerChan := make(chan string, workers)

	// Start workers
	for i := 0; i < workers; i++ {
		go func() {
			for target := range workerChan {
				result := d.Scan(target, port)
				resultsChan <- result
			}
		}()
	}

	// Send work
	go func() {
		for _, target := range targets {
			workerChan <- target
		}
		close(workerChan)
	}()

	// Collect results
	var results []ReconInfo
	for i := 0; i < len(targets); i++ {
		results = append(results, <-resultsChan)
	}

	close(resultsChan)
	return results
}

// buildX224ConnectionRequest builds X.224 Connection Request packet
func buildX224ConnectionRequest() []byte {
	// X.224 Connection Request PDU
	// This is a basic implementation
	pdu := []byte{
		0x03, 0x00, 0x00, 0x13, 0x0e, 0xe0, 0x00, 0x00,
		0x00, 0x00, 0x08, 0x02, 0xf0, 0x80, 0x7f, 0x65,
		0x82, 0x84, 0xe0,
	}
	return pdu
}

// detectNLAFromResponse detects NLA from X.224 response
func detectNLAFromResponse(data []byte) bool {
	// Look for NLA indicators in response
	// PROTOCOL_HYBRID flag indicates NLA support
	if len(data) > 10 {
		// Check for specific bytes indicating NLA
		for i := 0; i < len(data)-1; i++ {
			if data[i] == 0x02 && data[i+1] == 0x00 {
				return true
			}
		}
	}
	return false
}
