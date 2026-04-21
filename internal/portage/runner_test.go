package portage

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNewRunner_ReturnsRunner(t *testing.T) {
	tr := New(nil)
	r := NewRunner(tr, func() (*snapshot.Snapshot, error) { return nil, nil }, 0)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
	if r.interval != defaultInterval {
		t.Errorf("expected default interval %v, got %v", defaultInterval, r.interval)
	}
}

func TestRunner_CancelsCleanly(t *testing.T) {
	tr := New(nil)
	r := NewRunner(tr, func() (*snapshot.Snapshot, error) {
		return snapshot.New(nil), nil
	}, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not cancel within deadline")
	}
}

func TestRunner_InvokesSnapshotFunc(t *testing.T) {
	tr := New(nil)
	var calls atomic.Int32

	snap := func() (*snapshot.Snapshot, error) {
		calls.Add(1)
		return snapshot.New(nil), nil
	}

	r := NewRunner(tr, snap, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if calls.Load() < 2 {
		t.Errorf("expected at least 2 snapshot calls, got %d", calls.Load())
	}
}

func TestRunner_SkipsOnSnapshotError(t *testing.T) {
	tr := New(nil)
	var calls atomic.Int32

	snap := func() (*snapshot.Snapshot, error) {
		calls.Add(1)
		return nil, errors.New("scan failed")
	}

	r := NewRunner(tr, snap, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 65*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	// Tracker should have no data since every snapshot errored.
	if tr.Age(0, "tcp") != 0 {
		t.Error("expected zero age when all snapshots errored")
	}
	if calls.Load() < 2 {
		t.Errorf("expected at least 2 calls, got %d", calls.Load())
	}
}
