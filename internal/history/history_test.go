package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNew_NoFile(t *testing.T) {
	h, err := history.New(tempPath(t), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := h.Events(); len(got) != 0 {
		t.Fatalf("expected 0 events, got %d", len(got))
	}
}

func TestRecord_AppendsEvent(t *testing.T) {
	h, _ := history.New(tempPath(t), 10)
	if err := h.Record("opened", "tcp", 8080); err != nil {
		t.Fatalf("Record: %v", err)
	}
	events := h.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Port != 8080 || events[0].Kind != "opened" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestRecord_CapsAtMaxEvents(t *testing.T) {
	h, _ := history.New(tempPath(t), 3)
	for i := 0; i < 5; i++ {
		_ = h.Record("opened", "tcp", 8000+i)
	}
	if got := len(h.Events()); got != 3 {
		t.Fatalf("expected 3 events (cap), got %d", got)
	}
}

func TestRecord_KeepsNewest(t *testing.T) {
	h, _ := history.New(tempPath(t), 2)
	_ = h.Record("opened", "tcp", 1)
	_ = h.Record("opened", "tcp", 2)
	_ = h.Record("opened", "tcp", 3)
	events := h.Events()
	if events[0].Port != 2 || events[1].Port != 3 {
		t.Errorf("expected ports 2,3; got %d,%d", events[0].Port, events[1].Port)
	}
}

func TestNew_LoadsExisting(t *testing.T) {
	p := tempPath(t)
	h1, _ := history.New(p, 10)
	_ = h1.Record("closed", "udp", 53)

	h2, err := history.New(p, 10)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if got := len(h2.Events()); got != 1 {
		t.Fatalf("expected 1 persisted event, got %d", got)
	}
}

func TestNew_BadPath(t *testing.T) {
	_, err := history.New("/no/such/dir/h.json", 10)
	if err == nil {
		t.Fatal("expected error for unwritable path")
	}
}

func TestRecord_Persists(t *testing.T) {
	p := tempPath(t)
	h, _ := history.New(p, 10)
	_ = h.Record("opened", "tcp", 443)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
