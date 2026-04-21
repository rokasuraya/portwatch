package portflap

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

func TestNew_DefaultsToStderr(t *testing.T) {
	d := New(3, time.Minute)
	if d == nil {
		t.Fatal("expected non-nil Detector")
	}
	if d.Threshold != 3 {
		t.Errorf("threshold: got %d, want 3", d.Threshold)
	}
	if d.Window != time.Minute {
		t.Errorf("window: got %v, want 1m", d.Window)
	}
}

func TestObserve_NoWarnBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	d := New(4, time.Minute)
	d.SetOutput(&buf)

	e := makeEntry(8080, "tcp")
	d.Observe([]snapshot.Entry{e}, nil)
	d.Observe(nil, []snapshot.Entry{e})
	d.Observe([]snapshot.Entry{e}, nil)

	if buf.Len() != 0 {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestObserve_WarnsAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	d := New(3, time.Minute)
	d.SetOutput(&buf)

	e := makeEntry(22, "tcp")
	for i := 0; i < 3; i++ {
		d.Observe([]snapshot.Entry{e}, nil)
	}

	if !strings.Contains(buf.String(), "22/tcp") {
		t.Errorf("expected warning for 22/tcp, got: %s", buf.String())
	}
}

func TestObserve_PrunesOldEvents(t *testing.T) {
	var buf bytes.Buffer
	d := New(3, 10*time.Second)
	d.SetOutput(&buf)

	past := time.Now().Add(-20 * time.Second)
	e := makeEntry(443, "tcp")
	k := portKey(e)

	// Seed two old events manually.
	d.mu.Lock()
	d.counts[k] = []time.Time{past, past}
	d.mu.Unlock()

	// One new event — total within window should be 1, below threshold.
	d.Observe([]snapshot.Entry{e}, nil)

	if buf.Len() != 0 {
		t.Errorf("old events should have been pruned; got: %s", buf.String())
	}
}

func TestObserve_ClosedCountsAsTransition(t *testing.T) {
	var buf bytes.Buffer
	d := New(2, time.Minute)
	d.SetOutput(&buf)

	e := makeEntry(3306, "tcp")
	d.Observe([]snapshot.Entry{e}, nil)
	d.Observe(nil, []snapshot.Entry{e})

	if !strings.Contains(buf.String(), "3306/tcp") {
		t.Errorf("expected warning for 3306/tcp, got: %s", buf.String())
	}
}

func TestReset_ClearsState(t *testing.T) {
	var buf bytes.Buffer
	d := New(2, time.Minute)
	d.SetOutput(&buf)

	e := makeEntry(9200, "tcp")
	d.Observe([]snapshot.Entry{e}, nil)
	d.Observe([]snapshot.Entry{e}, nil) // would warn
	buf.Reset()
	d.Reset()

	// After reset a single event must not warn.
	d.Observe([]snapshot.Entry{e}, nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output after reset; got: %s", buf.String())
	}
}
