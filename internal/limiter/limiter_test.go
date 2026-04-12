package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/limiter"
)

func TestNew_DefaultsToOne(t *testing.T) {
	l := limiter.New(0)
	if l.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", l.Cap())
	}
}

func TestNew_SetsCapCorrectly(t *testing.T) {
	l := limiter.New(5)
	if l.Cap() != 5 {
		t.Fatalf("expected cap 5, got %d", l.Cap())
	}
}

func TestAcquireRelease_SingleSlot(t *testing.T) {
	l := limiter.New(1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Available() != 0 {
		t.Fatalf("expected 0 available after acquire, got %d", l.Available())
	}
	l.Release()
	if l.Available() != 1 {
		t.Fatalf("expected 1 available after release, got %d", l.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	l := limiter.New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Fill the single slot.
	_ = l.Acquire(context.Background())

	// Second acquire should block and return ctx.Err().
	err := l.Acquire(ctx)
	if err == nil {
		t.Fatal("expected error when context cancelled, got nil")
	}
}

func TestAcquire_CancelledContextReturnsError(t *testing.T) {
	l := limiter.New(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := l.Acquire(ctx)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestLimiter_ConcurrencyRespected(t *testing.T) {
	const cap = 3
	l := limiter.New(cap)
	ctx := context.Background()

	var (
		mu      sync.Mutex
		peak    int
		current int
		wg      sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(ctx)
			defer l.Release()

			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			mu.Lock()
			current--
			mu.Unlock()
		}()
	}
	wg.Wait()

	if peak > cap {
		t.Fatalf("peak concurrency %d exceeded cap %d", peak, cap)
	}
}
