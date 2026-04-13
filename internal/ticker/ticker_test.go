package ticker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/ticker"
)

func TestNew_ReturnsTicker(t *testing.T) {
	tk := ticker.New(time.Second, 0, func(ctx context.Context) {})
	if tk == nil {
		t.Fatal("expected non-nil Ticker")
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	tk := ticker.New(10*time.Millisecond, 0, func(ctx context.Context) {})
	go func() {
		tk.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestRun_InvokesOnTick(t *testing.T) {
	var count atomic.Int64
	ctx, cancel := context.WithCancel(context.Background())

	tk := ticker.New(20*time.Millisecond, 0, func(ctx context.Context) {
		count.Add(1)
	})

	go tk.Run(ctx)
	time.Sleep(90 * time.Millisecond)
	cancel()

	got := count.Load()
	if got < 2 {
		t.Fatalf("expected at least 2 ticks, got %d", got)
	}
}

func TestNew_JitterClampedWhenEqualToInterval(t *testing.T) {
	// Should not panic; jitter >= interval is clamped to zero.
	var fired atomic.Bool
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	tk := ticker.New(30*time.Millisecond, 30*time.Millisecond, func(ctx context.Context) {
		fired.Store(true)
	})
	go tk.Run(ctx)
	<-ctx.Done()

	if !fired.Load() {
		t.Fatal("expected at least one tick to fire")
	}
}

func TestRun_JitterDoesNotExceedInterval(t *testing.T) {
	// Verify the ticker still fires within a reasonable window when jitter is set.
	var count atomic.Int64
	ctx, cancel := context.WithCancel(context.Background())

	tk := ticker.New(30*time.Millisecond, 5*time.Millisecond, func(ctx context.Context) {
		count.Add(1)
	})

	go tk.Run(ctx)
	time.Sleep(120 * time.Millisecond)
	cancel()

	if count.Load() < 2 {
		t.Fatalf("expected ticks with jitter, got %d", count.Load())
	}
}
