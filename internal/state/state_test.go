package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestNew_NoExistingFile(t *testing.T) {
	s, err := New(tempFile(t))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(s.Current().Ports) != 0 {
		t.Error("expected empty initial state")
	}
}

func TestUpdate_DetectsOpenedPorts(t *testing.T) {
	s, _ := New(tempFile(t))

	ports := []PortState{
		{Port: 80, Protocol: "tcp", Open: true, SeenAt: time.Now()},
		{Port: 443, Protocol: "tcp", Open: true, SeenAt: time.Now()},
	}
	diff, err := s.Update(ports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(diff.Opened))
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(diff.Closed))
	}
}

func TestUpdate_DetectsClosedPorts(t *testing.T) {
	s, _ := New(tempFile(t))

	initial := []PortState{
		{Port: 80, Protocol: "tcp", Open: true, SeenAt: time.Now()},
		{Port: 8080, Protocol: "tcp", Open: true, SeenAt: time.Now()},
	}
	s.Update(initial) //nolint

	updated := []PortState{
		{Port: 80, Protocol: "tcp", Open: true, SeenAt: time.Now()},
	}
	diff, err := s.Update(updated)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Closed) != 1 || diff.Closed[0].Port != 8080 {
		t.Errorf("expected port 8080 to be closed, got %+v", diff.Closed)
	}
}

func TestUpdate_PersistsState(t *testing.T) {
	path := tempFile(t)
	s, _ := New(path)

	ports := []PortState{
		{Port: 22, Protocol: "tcp", Open: true, SeenAt: time.Now()},
	}
	s.Update(ports) //nolint

	// Reload from disk
	s2, err := New(path)
	if err != nil {
		t.Fatalf("failed to reload state: %v", err)
	}
	if _, ok := s2.Current().Ports[22]; !ok {
		t.Error("expected port 22 to be persisted")
	}
}

func TestNew_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "bad.json")
	os.WriteFile(badPath, []byte("not json{"), 0644)

	_, err := New(badPath)
	if err == nil {
		t.Error("expected error for invalid JSON file")
	}
}
