package processor

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Processor struct {
	dirInput  string
	dirOutput string
	dirLog    string
	reports   map[string][]Report
}

func NewProcessor(dirInput, dirOutput, dirLog string) *Processor {
	return &Processor{
		dirInput:  dirInput,
		dirOutput: dirOutput,
		dirLog:    dirLog,
		reports:   make(map[string][]Report),
	}
}

// ProcessZipFiles -> processZipFile -> extractAndCount -> aggregateData
func (p *Processor) ProcessZipFiles() error {
	files, err := os.ReadDir(p.dirInput)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".zip" {
			if err := p.processZipFile(filepath.Join(p.dirInput, file.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) processZipFile(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		// creates a new report for the internal file
		report := NewReport(f.Name)

		// fill the report
		if err := p.extractAndCount(f, report); err != nil {
			fmt.Printf("failed to process %s -> %s: %v\n", zipPath, f.Name, err)
		}

		// add report to the reports map
		p.reports[zipPath] = append(p.reports[zipPath], *report)
	}

	return nil
}

func (p *Processor) extractAndCount(f *zip.File, report *Report) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	if strings.HasSuffix(f.Name, ".csv") {
		return nil
	}

	reader := csv.NewReader(rc)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	lineCount := len(records) - 1 // Ignore header
	report.AddLineCount(f.Name, lineCount)

	if err := p.aggregateData(f.Name, records[1:]); err != nil {
		return err
	}

	return nil
}

func (p *Processor) aggregateData(fileName string, records [][]string) error {
	outputFilePath := filepath.Join(p.dirOutput, fileName)
	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// GenerateReport generates and saves the final report to the log directory.
func (p *Processor) GenerateReport() error {

	var sb strings.Builder

	// Header
	currentTime := time.Now().Format("2006-01-02 15:04:05.000000")
	sb.WriteString("# REPORT\n")
	sb.WriteString(fmt.Sprintf("Generated at: %s\n", currentTime))
	sb.WriteString(fmt.Sprintf("DIR: %s\n", p.dirInput))

	// Body
	for zipFileName, reportList := range p.reports {
		sb.WriteString(fmt.Sprintf("\n\n## ZIP: %s\n", filepath.Base(zipFileName)))

		for _, report := range reportList {
			fmt.Println()
			sb.WriteString(report.String())
		}
	}

	// Save the reports to the log directory
	reportFilePath := filepath.Join(p.dirLog, "report.txt")
	file, err := os.OpenFile(reportFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create report file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(sb.String())
	if err != nil {
		return fmt.Errorf("failed to write to report file: %w", err)
	}

	// Mostra na tela
	fmt.Println(sb.String())

	return nil
}
