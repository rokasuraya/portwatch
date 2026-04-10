package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Report represents a single scan report entry.
type Report struct {
	Timestamp time.Time `json:"timestamp"`
	OpenedPorts []int     `json:"opened_ports"`
	ClosedPorts []int     `json:"closed_ports"`
	TotalOpen   int       `json:"total_open"`
}

// Reporter writes scan reports to a destination.
type Reporter struct {
	writer io.Writer
	jsonMode bool
}

// New creates a Reporter. If outputPath is empty, stdout is used.
// If jsonMode is true, reports are written as JSON lines.
func New(outputPath string, jsonMode bool) (*Reporter, error) {
	var w io.Writer = os.Stdout
	if outputPath != "" {
		f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("reporter: open output file: %w", err)
		}
		w = f
	}
	return &Reporter{writer: w, jsonMode: jsonMode}, nil
}

// Write emits a report to the configured destination.
func (r *Reporter) Write(report Report) error {
	if r.jsonMode {
		return r.writeJSON(report)
	}
	return r.writeText(report)
}

func (r *Reporter) writeJSON(report Report) error {
	data, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("reporter: marshal report: %w", err)
	}
	_, err = fmt.Fprintf(r.writer, "%s\n", data)
	return err
}

func (r *Reporter) writeText(report Report) error {
	_, err := fmt.Fprintf(
		r.writer,
		"[%s] opened=%d closed=%d total_open=%d\n",
		report.Timestamp.Format(time.RFC3339),
		len(report.OpenedPorts),
		len(report.ClosedPorts),
		report.TotalOpen,
	)
	return err
}
