package processor

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	report := NewReport("text_report.txt")

	report.AddFileReport("email_list.zip")

	report.AddLineCount("email_list.zip", 200)
	report.LogError("error 1")
	report.LogError("error 2")

	output := report.String()

	if len(report.LineCounts) == 0 {
		t.Errorf("expected LineCounts to be 200, got %d items", len(report.LineCounts))
	}
	if len(report.Errors) == 0 {
		t.Errorf("expected Errors to be 2, got %d items", len(report.Errors))
	}

	if !strings.Contains(output, "├── File: email_list.zip: 200 lines") {
		t.Errorf("expected output to contain file info, got: %s", output)
	}
	if !strings.Contains(output, "│    └── Errors:") {
		t.Errorf("expected output to contain errors section, got: %s", output)
	}
	if !strings.Contains(output, "│         Errors  error 1") || !strings.Contains(output, "│         Errors  error 2") {
		t.Errorf("expected output to contain specific errors, got: %s", output)
	}
}
