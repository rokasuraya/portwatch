package reporter

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	r, err := New("", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
	if r.writer != os.Stdout {
		t.Error("expected writer to be stdout")
	}
}

func TestNew_CreatesOutputFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.log")
	r, err := New(path, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestWrite_TextMode(t *testing.T) {
	var buf bytes.Buffer
	r := &Reporter{writer: &buf, jsonMode: false}
	report := Report{
		Timestamp:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		OpenedPorts: []int{8080},
		ClosedPorts: []int{},
		TotalOpen:   5,
	}
	if err := r.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "opened=1") {
		t.Errorf("expected opened=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "total_open=5") {
		t.Errorf("expected total_open=5 in output, got: %s", out)
	}
}

func TestWrite_JSONMode(t *testing.T) {
	var buf bytes.Buffer
	r := &Reporter{writer: &buf, jsonMode: true}
	report := Report{
		Timestamp:   time.Now(),
		OpenedPorts: []int{22, 80},
		ClosedPorts: []int{443},
		TotalOpen:   10,
	}
	if err := r.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Report
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &decoded); err != nil {
		t.Fatalf("failed to decode JSON output: %v", err)
	}
	if decoded.TotalOpen != 10 {
		t.Errorf("expected TotalOpen=10, got %d", decoded.TotalOpen)
	}
	if len(decoded.OpenedPorts) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(decoded.OpenedPorts))
	}
}
