package credentials

import (
	"falcon/config"
	"fmt"
	"os"
	"strings"
)

// GenerateCredentials creates credential combinations from users and passwords
func GenerateCredentials(users []string, passwords []string, domain string) []*config.Credential {
	credentials := make([]*config.Credential, 0)

	for _, user := range users {
		// Extract domain if present
		userDomain := domain
		if strings.Contains(user, "\\") {
			parts := strings.Split(user, "\\")
			userDomain = parts[0]
			user = parts[1]
		} else if strings.Contains(user, "@") {
			parts := strings.Split(user, "@")
			user = parts[0]
			userDomain = parts[1]
		}

		for _, pass := range passwords {
			cred := &config.Credential{
				Username: user,
				Password: pass,
				Domain:   userDomain,
			}
			credentials = append(credentials, cred)
		}
	}

	return credentials
}

// SaveCredentials saves generated credentials to file
func SaveCredentials(credentials []*config.Credential, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create credentials file: %w", err)
	}
	defer file.Close()

	for _, cred := range credentials {
		line := fmt.Sprintf("%s:%s", cred.Username, cred.Password)
		if cred.Domain != "" {
			line = fmt.Sprintf("%s\\\\%s", cred.Domain, line)
		}
		file.WriteString(line + "\n")
	}

	return nil
}

// LoadCredentials loads credentials from file
func LoadCredentials(filename string) ([]*config.Credential, error) {
	credentials := make([]*config.Credential, 0)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		cred := parseCredentialLine(line)
		if cred != nil {
			credentials = append(credentials, cred)
		}
	}

	return credentials, nil
}

func parseCredentialLine(line string) *config.Credential {
	domain := ""

	// Check for domain\user:pass format
	if strings.Contains(line, "\\") {
		parts := strings.Split(line, "\\")
		if len(parts) == 2 {
			domain = parts[0]
			line = parts[1]
		}
	}

	// Parse user:pass
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return nil
	}

	return &config.Credential{
		Username: parts[0],
		Password: parts[1],
		Domain:   domain,
	}
}
