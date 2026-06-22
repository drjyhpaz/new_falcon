package credentials

import (
	"bufio"
	"falcon/config"
	"fmt"
	"os"
	"strings"
)

// LoadServers loads targets from servers.txt
func LoadServers(filename string) ([]*config.Target, error) {
	var targets []*config.Target

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open servers file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		ip := strings.TrimSpace(parts[0])
		var port int
		fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &port)

		if port <= 0 || port > 65535 {
			continue
		}

		target := &config.Target{
			IP:   ip,
			Port: port,
		}
		targets = append(targets, target)
	}

	return targets, scanner.Err()
}

// LoadUsers loads usernames from users.txt
func LoadUsers(filename string) ([]string, error) {
	var users []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open users file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			users = append(users, line)
		}
	}

	return users, scanner.Err()
}

// LoadPasswords loads passwords from passwords.txt
func LoadPasswords(filename string) ([]string, error) {
	var passwords []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open passwords file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			passwords = append(passwords, line)
		}
	}

	return passwords, scanner.Err()
}
