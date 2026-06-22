package utils

import (
	"fmt"
	"net"
	"strings"
)

// ValidateIP validates an IP address
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ValidatePort validates a port number
func ValidatePort(port int) bool {
	return port > 0 && port <= 65535
}

// ValidateCIDR validates a CIDR notation
func ValidateCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// ExpandCIDR expands a CIDR notation to individual IPs
func ExpandCIDR(cidr string) ([]string, error) {
	var ips []string

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ExtractDomain extracts domain from username
func ExtractDomain(username string) (string, string) {
	if strings.Contains(username, "\\") {
		parts := strings.Split(username, "\\")
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}

	if strings.Contains(username, "@") {
		parts := strings.Split(username, "@")
		if len(parts) == 2 {
			return parts[1], parts[0]
		}
	}

	return "", username
}

// FormatTarget formats a target for display
func FormatTarget(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

// FormatCredential formats a credential for display
func FormatCredential(username, password, domain string) string {
	if domain != "" {
		return fmt.Sprintf("%s\\\\%s:%s", domain, username, password)
	}
	return fmt.Sprintf("%s:%s", username, password)
}
