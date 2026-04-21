package portcap_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"portwatch/internal/portcap"
	"portwatch/internal/snapshot"
)

func TestNewRunner_ReturnsRunner(t *testing.T) {
	c := portcap.New(10, nil)
	fn := func() (*snapshot.Snapshot, error) { return snapshot.New(nil), nil }
	r := portcap.NewRunner(c, fn, time.Second)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestRunner_CancelsCleanly(t *testing.T) {
	c := portcap.New(10, nil)
	fn := func() (*snapshot.Snapshot, error) { return snapshot.New(nil), nil }
	r := portcap.NewRunner(c, fn, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := r.Run(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestRunner_SkipsOnSnapshotError(t *testing.T) {
	c := portcap.New(10, nil)
	calls := 0
	fn := func() (*snapshot.Snapshot, error) {
		calls++
		return nil, errors.New("scan failed")
	}
	r := portcap.NewRunner(c, fn, 30*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = r.Run(ctx)
	if calls == 0 {
		t.Fatal("expected snapshot function to be called at least once")
	}
}

func TestRunner_DefaultIntervalApplied(t *testing.T) {
	c := portcap.New(10, nil)
	fn := func() (*snapshot.Snapshot, error) { return snapshot.New(nil), nil }
	// zero interval should default internally without panic
	r := portcap.NewRunner(c, fn, 0)
	if r == nil {
		t.Fatal("expected runner with default interval")
	}
}
