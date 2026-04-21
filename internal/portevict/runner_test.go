package portevict

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNewRunner_ReturnsRunner(t *testing.T) {
	ev := New(time.Minute)
	r := NewRunner(ev, func() (*snapshot.Snapshot, error) { return nil, nil }, time.Second)
	if r == nil {
		t.Fatal("expected non-nil Runner")
	}
}

func TestRunner_CancelsCleanly(t *testing.T) {
	ev := New(time.Minute)
	r := NewRunner(ev, func() (*snapshot.Snapshot, error) { return nil, nil }, 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not cancel in time")
	}
}

func TestRunner_SkipsOnSnapshotError(t *testing.T) {
	ev := New(time.Minute)
	called := 0
	r := NewRunner(ev, func() (*snapshot.Snapshot, error) {
		called++
		return nil, errors.New("boom")
	}, 10*time.Millisecond)
	r.OnReturn = func(snapshot.Entry) { t.Error("OnReturn should not be called on error") }
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()
	r.Run(ctx)
	if called == 0 {
		t.Fatal("expected snapshot func to be called")
	}
}
