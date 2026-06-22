package credentials

import (
	"fmt"
	"strings"
	"sync"

	"github.com/falconjonz/falcon_rdp/config"
	"github.com/falconjonz/falcon_rdp/utils"
)

type CredentialLoader struct {
	users       []string
	passwords   []string
	credentials []config.Credential
	mu          sync.RWMutex
}

// NewCredentialLoader creates a new credential loader
func NewCredentialLoader() *CredentialLoader {
	return &CredentialLoader{}
}

// LoadUsers loads usernames from file
func (cl *CredentialLoader) LoadUsers(filename string) error {
	lines, err := utils.ReadFile(filename)
	if err != nil {
		return err
	}

	cl.mu.Lock()
	cl.users = lines
	cl.mu.Unlock()

	return nil
}

// LoadPasswords loads passwords from file
func (cl *CredentialLoader) LoadPasswords(filename string) error {
	lines, err := utils.ReadFile(filename)
	if err != nil {
		return err
	}

	cl.mu.Lock()
	cl.passwords = lines
	cl.mu.Unlock()

	return nil
}

// LoadCredentials loads pre-generated credentials from file
func (cl *CredentialLoader) LoadCredentials(filename string) error {
	lines, err := utils.ReadFile(filename)
	if err != nil {
		return err
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()

	for _, line := range lines {
		domain, user := utils.ParseDomain(line)
		parts := strings.SplitN(user, ":", 2)
		if len(parts) == 2 {
			cl.credentials = append(cl.credentials, config.Credential{
				Username: parts[0],
				Password: parts[1],
				Domain:   domain,
			})
		}
	}

	return nil
}

// GenerateCredentials generates cartesian product of users and passwords
func (cl *CredentialLoader) GenerateCredentials(defaultDomain string) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if len(cl.users) == 0 || len(cl.passwords) == 0 {
		return fmt.Errorf("users or passwords not loaded")
	}

	cl.credentials = make([]config.Credential, 0, len(cl.users)*len(cl.passwords))

	for _, user := range cl.users {
		domain, username := utils.ParseDomain(user)
		if domain == "" {
			domain = defaultDomain
		}

		for _, pass := range cl.passwords {
			cl.credentials = append(cl.credentials, config.Credential{
				Username: username,
				Password: pass,
				Domain:   domain,
			})
		}
	}

	return nil
}

// GetCredentials returns all credentials
func (cl *CredentialLoader) GetCredentials() []config.Credential {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	// Return a copy
	creds := make([]config.Credential, len(cl.credentials))
	copy(creds, cl.credentials)
	return creds
}

// GetUsers returns all users
func (cl *CredentialLoader) GetUsers() []string {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	users := make([]string, len(cl.users))
	copy(users, cl.users)
	return users
}

// GetPasswords returns all passwords
func (cl *CredentialLoader) GetPasswords() []string {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	passes := make([]string, len(cl.passwords))
	copy(passes, cl.passwords)
	return passes
}

// Count returns credential count
func (cl *CredentialLoader) Count() int {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	return len(cl.credentials)
}
