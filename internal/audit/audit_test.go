package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	a := audit.New(nil)
	if a == nil {
		t.Fatal("expected non-nil Audit")
	}
}

func TestLog_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)

	if err := a.Log("opened", "tcp", 8080, "test"); err != nil {
		t.Fatalf("Log returned error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var entry audit.Entry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if entry.Event != "opened" {
		t.Errorf("expected event=opened, got %q", entry.Event)
	}
	if entry.Port != 8080 {
		t.Errorf("expected port=8080, got %d", entry.Port)
	}
	if entry.Protocol != "tcp" {
		t.Errorf("expected protocol=tcp, got %q", entry.Protocol)
	}
	if entry.Note != "test" {
		t.Errorf("expected note=test, got %q", entry.Note)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogOpened_SetsEvent(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	_ = a.LogOpened("udp", 53)

	var entry audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry)
	if entry.Event != "opened" {
		t.Errorf("expected opened, got %q", entry.Event)
	}
}

func TestLogClosed_SetsEvent(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	_ = a.LogClosed("tcp", 443)

	var entry audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry)
	if entry.Event != "closed" {
		t.Errorf("expected closed, got %q", entry.Event)
	}
}

func TestLog_MultipleEntriesOnePerLine(t *testing.T) {
	var buf bytes.Buffer
	a := audit.New(&buf)
	_ = a.LogOpened("tcp", 80)
	_ = a.LogClosed("tcp", 22)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}
