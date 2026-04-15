package stale

import (
	"bytes"
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNewRunner_ReturnsRunner(t *testing.T) {
	d := New(time.Minute, &bytes.Buffer{})
	fn := func(_ context.Context) (*snapshot.Snapshot, error) { return snapshot.New(nil), nil }
	r := NewRunner(d, fn, time.Second)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestRunner_CancelsCleanly(t *testing.T) {
	d := New(time.Minute, &bytes.Buffer{})
	fn := func(_ context.Context) (*snapshot.Snapshot, error) { return snapshot.New(nil), nil }
	r := NewRunner(d, fn, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not cancel within timeout")
	}
}

func TestRunner_InvokesSnapshotFunc(t *testing.T) {
	var calls atomic.Int32
	buf := &bytes.Buffer{}
	d := New(time.Minute, buf)

	fn := func(_ context.Context) (*snapshot.Snapshot, error) {
		calls.Add(1)
		return snapshot.New([]snapshot.Entry{{Proto: "tcp", Port: 80}}), nil
	}

	r := NewRunner(d, fn, 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if calls.Load() < 2 {
		t.Fatalf("expected at least 2 snapshot calls, got %d", calls.Load())
	}
}

func TestRunner_SkipsOnSnapshotError(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New(time.Minute, buf)

	errFn := func(_ context.Context) (*snapshot.Snapshot, error) {
		return nil, errors.New("scan failed")
	}

	r := NewRunner(d, errFn, 10*time.Millisecond)
	ctx, cancel := context.With35*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	if len(d.seen) != 0 {
		t.Fatal("expected no ports tracked after error")
	}
}
