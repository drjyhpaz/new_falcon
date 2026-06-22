package report

import (
	"encoding/csv"
	"encoding/json"
	"falcon/config"
	"fmt"
	"os"
	"time"
)

// ReportGenerator generates attack reports
type ReportGenerator struct {
	Results    []*config.Result
	Statistics map[string]interface{}
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(results []*config.Result) *ReportGenerator {
	return &ReportGenerator{
		Results:    results,
		Statistics: make(map[string]interface{}),
	}
}

// GenerateJSONReport generates a JSON report
func (rg *ReportGenerator) GenerateJSONReport(filename string) error {
	report := map[string]interface{}{
		"generated_at": time.Now(),
		"total_results": len(rg.Results),
		"successful_logins": rg.countSuccessful(),
		"results": rg.Results,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// GenerateCSVReport generates a CSV report
func (rg *ReportGenerator) GenerateCSVReport(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"IP", "Port", "Username", "Password", "Domain", "Status", "Timestamp", "Error"}
	writer.Write(header)

	// Write results
	for _, result := range rg.Results {
		status := "FAILED"
		if result.Success {
			status = "SUCCESS"
		}

		row := []string{
			result.IP,
			fmt.Sprintf("%d", result.Port),
			result.Username,
			result.Password,
			result.Domain,
			status,
			result.Timestamp.Format(time.RFC3339),
			result.Error,
		}
		writer.Write(row)
	}

	return nil
}

// GenerateSummary generates a summary report
func (rg *ReportGenerator) GenerateSummary() string {
	successful := rg.countSuccessful()
	total := len(rg.Results)

	summary := fmt.Sprintf(`
╔════════════════════════════════════════════╗
║       FALCON ATTACK REPORT SUMMARY          ║
╠════════════════════════════════════════════╣
║ Total Attempts:        %-26d ║
║ Successful Logins:     %-26d ║
║ Success Rate:          %.2f%%                 ║
║ Generated:             %-26s ║
╚════════════════════════════════════════════╝
`, total, successful, float64(successful)/float64(total)*100, time.Now().Format("2006-01-02 15:04:05"))

	return summary
}

func (rg *ReportGenerator) countSuccessful() int {
	count := 0
	for _, result := range rg.Results {
		if result.Success {
			count++
		}
	}
	return count
}
