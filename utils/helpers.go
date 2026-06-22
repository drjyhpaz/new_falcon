package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// ParseTargetString parses IP:Port format
func ParseTargetString(target string) (string, uint16, error) {
	parts := strings.Split(strings.TrimSpace(target), ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid target format: %s", target)
	}

	ip := parts[0]
	if net.ParseIP(ip) == nil {
		return "", 0, fmt.Errorf("invalid IP: %s", ip)
	}

	var port uint16
	_, err := fmt.Sscanf(parts[1], "%d", &port)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %s", parts[1])
	}

	return ip, port, nil
}

// ParseDomain extracts domain from credential (DOMAIN\user or user@domain)
func ParseDomain(credential string) (string, string) {
	if strings.Contains(credential, "\\") {
		parts := strings.Split(credential, "\\")
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}

	if strings.Contains(credential, "@") {
		parts := strings.Split(credential, "@")
		if len(parts) == 2 {
			return parts[1], parts[0]
		}
	}

	return "", credential
}

// ReadFile reads lines from a file
func ReadFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

// WriteFile writes data to file (JSON format)
func WriteFile(filename string, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

// FileExists checks if file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// IsValidIP checks if string is valid IP
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// FormatBytes formats bytes to human readable
func FormatBytes(bytes int64) string {
	suffixes := []string{"B", "KB", "MB", "GB"}
	float := float64(bytes)
	for i := 0; i < len(suffixes)-1; i++ {
		if float < 1024 {
			return fmt.Sprintf("%.2f %s", float, suffixes[i])
		}
		float /= 1024
	}
	return fmt.Sprintf("%.2f %s", float, suffixes[len(suffixes)-1])
}

// GetPPS calculates packets per second
func GetPPS(attempts int64, duration time.Duration) float64 {
	if duration.Seconds() == 0 {
		return 0
	}
	return float64(attempts) / duration.Seconds()
}

// TimeSinceStart returns formatted time since start
func TimeSinceStart(start time.Time) string {
	duration := time.Since(start)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
