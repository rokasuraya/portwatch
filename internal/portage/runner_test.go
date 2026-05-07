package portage

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/snapshot"
)

func TestNewRunner_ReturnsRunner(t *testing.T) {
	tr := New(nil)
	r := NewRunner(tr, func() (*snapshot.Snapshot, error) { return nil, nil }, 0)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
	if r.interval != 30*time.Second {
		t.Errorf("default interval = %v, want 30s", r.interval)
	}
}

func TestRunner_CancelsCleanly(t *testing.T) {
	tr := New(fixedNow(epoch))
	r := NewRunner(tr, func() (*snapshot.Snapshot, error) {
		return makeSnap([]int{80}, "tcp"), nil
	}, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
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
		t.Fatal("Run did not cancel in time")
	}
}

func TestRunner_InvokesSnapshotFunc(t *testing.T) {
	var calls atomic.Int32
	tr := New(fixedNow(epoch))
	r := NewRunner(tr, func() (*snapshot.Snapshot, error) {
		calls.Add(1)
		return makeSnap([]int{443}, "tcp"), nil
	}, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if calls.Load() == 0 {
		t.Error("expected snapshot func to be called at least once")
	}
}

func TestRunner_SkipsOnSnapshotError(t *testing.T) {
	tr := New(fixedNow(epoch))
	r := NewRunner(tr, func() (*snapshot.Snapshot, error) {
		return nil, errors.New("scan failed")
	}, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if len(tr.All()) != 0 {
		t.Error("expected tracker to remain empty on repeated errors")
	}
}
