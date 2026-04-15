package grace_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/grace"
)

// TestIntegration_ConcurrentWorkersDrainCleanly spawns several concurrent
// workers and verifies that Shutdown blocks until every one has released.
func TestIntegration_ConcurrentWorkersDrainCleanly(t *testing.T) {
	const n = 20
	g := grace.New(2 * time.Second)

	var completed atomic.Int32

	for i := 0; i < n; i++ {
		if !g.Acquire() {
			t.Fatalf("worker %d: Acquire returned false before shutdown", i)
		}
		go func() {
			defer g.Release()
			time.Sleep(30 * time.Millisecond)
			completed.Add(1)
		}()
	}

	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}

	if got := completed.Load(); got != n {
		t.Errorf("expected %d completions, got %d", n, got)
	}
}

// TestIntegration_NoWorkersShutdownImmediate verifies that Shutdown returns
// immediately when there are no in-flight workers.
func TestIntegration_NoWorkersShutdownImmediate(t *testing.T) {
	g := grace.New(time.Second)

	start := time.Now()
	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Errorf("shutdown took too long with no workers: %v", elapsed)
	}
}
