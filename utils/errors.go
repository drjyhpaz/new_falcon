package utils

import (
	"fmt"
	"strings"
)

// CustomError represents a custom error type
type CustomError struct {
	Code    int
	Message string
	Details string
}

// Error implements the error interface
func (ce *CustomError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", ce.Code, ce.Message, ce.Details)
}

// ErrorCode constants
const (
	ErrInvalidTarget = iota
	ErrConnectionFailed
	ErrAuthenticationFailed
	ErrTimeoutExceeded
	ErrProxyFailed
	ErrPortClosed
	ErrNLAEnabled
	ErrSSLCertError
)

// ErrorMessages maps error codes to messages
var ErrorMessages = map[int]string{
	ErrInvalidTarget:       "Invalid Target",
	ErrConnectionFailed:    "Connection Failed",
	ErrAuthenticationFailed: "Authentication Failed",
	ErrTimeoutExceeded:     "Timeout Exceeded",
	ErrProxyFailed:         "Proxy Failed",
	ErrPortClosed:          "Port Closed",
	ErrNLAEnabled:          "NLA Enabled",
	ErrSSLCertError:        "SSL Certificate Error",
}

// ClassifyError classifies an error
func ClassifyError(err error) int {
	errStr := err.Error()

	if strings.Contains(errStr, "connection refused") {
		return ErrPortClosed
	} else if strings.Contains(errStr, "timeout") {
		return ErrTimeoutExceeded
	} else if strings.Contains(errStr, "certificate") {
		return ErrSSLCertError
	}

	return ErrConnectionFailed
}

// NewCustomError creates a new custom error
func NewCustomError(code int, details string) *CustomError {
	message, ok := ErrorMessages[code]
	if !ok {
		message = "Unknown Error"
	}

	return &CustomError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
