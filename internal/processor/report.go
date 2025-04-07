package processor

import (
	"fmt"
	"strings"
)

// Report struct holds information about processed files and statistics.
type Report struct {
	FileName   string
	LineCounts map[string]int
	Errors     []string
}

// NewReport initializes a new Report instance.
func NewReport(fileName string) *Report {
	return &Report{
		FileName:   fileName,
		LineCounts: make(map[string]int),
		Errors:     []string{},
	}
}

// AddFileReport adds information about a processed ZIP file and its internal files to the report.
func (r *Report) AddFileReport(zipFileName string) {
	r.FileName = zipFileName
}

// AddLineCount adds the line count for a specific file type.
func (r *Report) AddLineCount(fileType string, count int) {
	r.LineCounts[fileType] = count
}

// LogError adds an error message to the report.
func (r *Report) LogError(message string) {
	r.Errors = append(r.Errors, message)
}

// String formats the report as a string for display and saving.
func (r *Report) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("├── File: %s: %d lines\n", r.FileName, r.LineCounts[r.FileName]))

	// Exibe os erros, se houver
	if len(r.Errors) > 0 {
		sb.WriteString("│    └── Errors:\n")
		for _, err := range r.Errors {
			sb.WriteString(fmt.Sprintf("│         Errors  %s\n", err))
		}
	}
	return sb.String()
}
