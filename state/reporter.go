package state

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/falconjonz/falcon_rdp/config"
)

// ReportGenerator generates attack reports
type ReportGenerator struct {
	results config.Result
	stats   config.Statistics
}

// GenerateJSONReport generates JSON format report
func GenerateJSONReport(filename string, results []config.Result, stats config.Statistics) error {
	report := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"statistics": map[string]interface{}{
			"total_attempts":      stats.TotalAttempts,
			"successful_logins":   stats.SuccessfulLogins,
			"failed_attempts":     stats.FailedAttempts,
			"duration":            stats.EndTime.Sub(stats.StartTime).Seconds(),
			"success_rate":        calculateSuccessRate(stats),
		},
		"results": results,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %v", err)
	}

	return nil
}

// GenerateCSVReport generates CSV format report
func GenerateCSVReport(filename string, results []config.Result) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"IP", "Port", "Username", "Password", "Domain", "Timestamp"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write results
	for _, result := range results {
		row := []string{
			result.IP,
			fmt.Sprintf("%d", result.Port),
			result.Username,
			result.Password,
			result.Domain,
			result.Timestamp.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	return nil
}

// GenerateTextReport generates plain text report
func GenerateTextReport(filename string, results []config.Result, stats config.Statistics) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Falcon RDP Brute-Force Report\n")
	fmt.Fprintf(file, "Generated: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "\n")

	// Statistics
	fmt.Fprintf(file, "=== STATISTICS ===\n")
	fmt.Fprintf(file, "Total Attempts: %d\n", stats.TotalAttempts)
	fmt.Fprintf(file, "Successful Logins: %d\n", stats.SuccessfulLogins)
	fmt.Fprintf(file, "Failed Attempts: %d\n", stats.FailedAttempts)
	fmt.Fprintf(file, "Success Rate: %.2f%%\n", calculateSuccessRate(stats))
	fmt.Fprintf(file, "Duration: %v\n", stats.EndTime.Sub(stats.StartTime))
	fmt.Fprintf(file, "\n")

	// Results
	if len(results) > 0 {
		fmt.Fprintf(file, "=== SUCCESSFUL LOGINS ===\n")
		for i, result := range results {
			fmt.Fprintf(file, "\n[%d] %s:%d\n", i+1, result.IP, result.Port)
			fmt.Fprintf(file, "    Username: %s\n", result.Username)
			fmt.Fprintf(file, "    Password: %s\n", result.Password)
			if result.Domain != "" {
				fmt.Fprintf(file, "    Domain: %s\n", result.Domain)
			}
			fmt.Fprintf(file, "    Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
		}
	} else {
		fmt.Fprintf(file, "=== NO SUCCESSFUL LOGINS ===\n")
	}

	return nil
}

// calculateSuccessRate calculates the success rate
func calculateSuccessRate(stats config.Statistics) float64 {
	if stats.TotalAttempts == 0 {
		return 0
	}
	return float64(stats.SuccessfulLogins) / float64(stats.TotalAttempts) * 100
}
