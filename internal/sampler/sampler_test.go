package sampler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/sampler"
)

// stubScanner satisfies sampler.Scanner.
type stubScanner struct {
	ports []string
	err   error
	calls int32
}

func (s *stubScanner) ScanPortRange(_, _ int) ([]string, error) {
	atomic.AddInt32(&s.calls, 1)
	return s.ports, s.err
}

func TestNew_ReturnsSampler(t *testing.T) {
	sc := &stubScanner{}
	sm := sampler.New(sc, 10*time.Millisecond, 0, nil)
	if sm == nil {
		t.Fatal("expected non-nil Sampler")
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	sc := &stubScanner{ports: []string{"tcp:80"}}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		sm := sampler.New(sc, 5*time.Millisecond, 0, nil)
		done <- sm.Run(ctx, 1, 1024)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestRun_InvokesOnSample(t *testing.T) {
	sc := &stubScanner{ports: []string{"tcp:443", "tcp:8080"}}

	var count int32
	onSample := func(ports []string) {
		atomic.AddInt32(&count, 1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sm := sampler.New(sc, 5*time.Millisecond, 0, onSample)
		_ = sm.Run(ctx, 1, 1024)
	}()

	time.Sleep(40 * time.Millisecond)
	cancel()

	if atomic.LoadInt32(&count) < 2 {
		t.Fatalf("expected at least 2 onSample calls, got %d", count)
	}
}

func TestRun_SkipsCallbackOnScanError(t *testing.T) {
	sc := &stubScanner{err: context.DeadlineExceeded}

	var count int32
	onSample := func(_ []string) { atomic.AddInt32(&count, 1) }

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	sm := sampler.New(sc, 5*time.Millisecond, 0, onSample)
	_ = sm.Run(ctx, 1, 1024)

	if atomic.LoadInt32(&count) != 0 {
		t.Fatalf("expected 0 onSample calls on error, got %d", count)
	}
}
