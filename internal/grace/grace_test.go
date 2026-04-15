package grace_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/grace"
)

func TestNew_DefaultTimeout(t *testing.T) {
	g := grace.New(0)
	if g == nil {
		t.Fatal("expected non-nil coordinator")
	}
}

func TestAcquireRelease_Nominal(t *testing.T) {
	g := grace.New(time.Second)
	if !g.Acquire() {
		t.Fatal("expected Acquire to return true before shutdown")
	}
	g.Release()

	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}
}

func TestAcquire_ReturnsFalseAfterShutdown(t *testing.T) {
	g := grace.New(time.Second)

	// start shutdown in background
	go g.Shutdown(context.Background()) //nolint:errcheck
	time.Sleep(20 * time.Millisecond)

	if g.Acquire() {
		t.Error("expected Acquire to return false after Shutdown called")
	}
}

func TestShutdown_WaitsForInFlight(t *testing.T) {
	g := grace.New(time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	g.Acquire()

	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		g.Release()
	}()

	start := time.Now()
	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Errorf("shutdown returned too quickly: %v", elapsed)
	}
	wg.Wait()
}

func TestShutdown_TimesOut(t *testing.T) {
	g := grace.New(50 * time.Millisecond)
	g.Acquire() // never released

	err := g.Shutdown(context.Background())
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestShutdown_ContextCancelled(t *testing.T) {
	g := grace.New(5 * time.Second)
	g.Acquire() // never released

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()

	err := g.Shutdown(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestShutdown_Idempotent(t *testing.T) {
	g := grace.New(time.Second)
	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
}
