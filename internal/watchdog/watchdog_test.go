package watchdog

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	wd := New(time.Second, nil)
	if wd.out == nil {
		t.Fatal("expected non-nil writer")
	}
	if wd.tolerance != time.Second {
		t.Fatalf("expected tolerance 1s, got %s", wd.tolerance)
	}
}

func TestBeat_UpdatesLastBeat(t *testing.T) {
	var buf bytes.Buffer
	wd := New(time.Second, &buf)

	before := time.Now()
	wd.Beat()
	after := time.Now()

	wd.mu.Lock()
	lb := wd.lastBeat
	wd.mu.Unlock()

	if lb.Before(before) || lb.After(after) {
		t.Fatalf("lastBeat %v not in [%v, %v]", lb, before, after)
	}
}

func TestCheck_NoWarnWhenFresh(t *testing.T) {
	var buf bytes.Buffer
	wd := New(500*time.Millisecond, &buf)
	wd.Beat()

	// check immediately — should be well within tolerance
	wd.check(time.Now(), 100*time.Millisecond)

	if buf.Len() != 0 {
		t.Fatalf("unexpected warning: %s", buf.String())
	}
}

func TestCheck_WarnsWhenStalled(t *testing.T) {
	var buf bytes.Buffer
	wd := New(0, &buf)

	// Manually set lastBeat to 5 seconds ago.
	wd.mu.Lock()
	wd.lastBeat = time.Now().Add(-5 * time.Second)
	wd.mu.Unlock()

	wd.check(time.Now(), 100*time.Millisecond)

	if !strings.Contains(buf.String(), "WARNING") {
		t.Fatalf("expected WARNING in output, got: %q", buf.String())
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	var buf bytes.Buffer
	wd := New(time.Second, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		wd.Run(ctx, 50*time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestRun_EmitsWarningOnMissedBeat(t *testing.T) {
	var buf bytes.Buffer
	wd := New(0, &buf)

	// Set lastBeat far in the past so the first check fires a warning.
	wd.mu.Lock()
	wd.lastBeat = time.Now().Add(-10 * time.Second)
	wd.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	wd.Run(ctx, 50*time.Millisecond)

	if !strings.Contains(buf.String(), "WARNING") {
		t.Fatalf("expected at least one WARNING, got: %q", buf.String())
	}
}
