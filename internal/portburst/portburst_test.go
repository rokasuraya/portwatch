package portburst

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port int, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func makeDiff(opened, closed []snapshot.Entry) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestNew_DefaultsToStderr(t *testing.T) {
	d := New(5, 30*time.Second, nil)
	if d.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestNew_DefaultThresholdAndWindow(t *testing.T) {
	d := New(0, 0, nil)
	if d.threshold != 5 {
		t.Fatalf("want threshold=5, got %d", d.threshold)
	}
	if d.window != 30*time.Second {
		t.Fatalf("want window=30s, got %v", d.window)
	}
}

func TestObserve_NoBurstBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	d := New(5, time.Minute, &buf)
	diff := makeDiff([]snapshot.Entry{
		makeEntry(8080, "tcp"),
		makeEntry(9090, "tcp"),
	}, nil)
	d.Observe(diff)
	if buf.Len() != 0 {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestObserve_WarnsAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	d := New(3, time.Minute, &buf)
	entries := []snapshot.Entry{
		makeEntry(80, "tcp"),
		makeEntry(443, "tcp"),
		makeEntry(22, "tcp"),
	}
	d.Observe(makeDiff(entries, nil))
	if !strings.Contains(buf.String(), "burst detected") {
		t.Fatalf("expected burst warning, got: %s", buf.String())
	}
}

func TestObserve_EmptyDiffNoOp(t *testing.T) {
	var buf bytes.Buffer
	d := New(2, time.Minute, &buf)
	d.Observe(makeDiff(nil, nil))
	if buf.Len() != 0 {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	var buf bytes.Buffer
	d := New(3, time.Minute, &buf)
	entries := []snapshot.Entry{
		makeEntry(80, "tcp"),
		makeEntry(443, "tcp"),
		makeEntry(22, "tcp"),
	}
	d.Observe(makeDiff(entries, nil))
	d.Reset()
	if d.Count(time.Now()) != 0 {
		t.Fatal("expected count=0 after reset")
	}
}

func TestCount_ExpiresOldEvents(t *testing.T) {
	var buf bytes.Buffer
	d := New(10, 5*time.Second, &buf)
	// Manually inject old events.
	old := time.Now().Add(-10 * time.Second)
	d.mu.Lock()
	d.events = []time.Time{old, old, old}
	d.mu.Unlock()
	if got := d.Count(time.Now()); got != 0 {
		t.Fatalf("want 0 active events, got %d", got)
	}
}
