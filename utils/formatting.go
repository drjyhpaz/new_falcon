package utils

import (
	"fmt"
	"strings"
	"time"
)

// FormatBytes converts bytes to human readable format
func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unitIndex := 0

	for size > 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

// FormatDuration converts duration to human readable format
func FormatDuration(d time.Duration) string {
	hours := d.Hours()
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", int(hours), minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// TruncateString truncates string to max length
func TruncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

// PadString pads string to specific width
func PadString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// RepeatString repeats string n times
func RepeatString(s string, n int) string {
	return strings.Repeat(s, n)
}

// ContainsAny checks if string contains any of the substrings
func ContainsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
