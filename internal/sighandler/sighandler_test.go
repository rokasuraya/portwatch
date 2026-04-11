package sighandler

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNew_DefaultSignals(t *testing.T) {
	h := New()
	if len(h.signals) != 2 {
		t.Fatalf("expected 2 default signals, got %d", len(h.signals))
	}
}

func TestNew_CustomSignals(t *testing.T) {
	h := New(syscall.SIGHUP)
	if len(h.signals) != 1 || h.signals[0] != syscall.SIGHUP {
		t.Fatalf("unexpected signals: %v", h.signals)
	}
}

func TestWait_CancelsOnContextDone(t *testing.T) {
	h := New()
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	go func() { done <- h.Wait(ctx) }()

	cancel()

	select {
	case sig := <-done:
		if sig != nil {
			t.Fatalf("expected nil signal on context cancel, got %v", sig)
		}
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after context cancel")
	}
}

func TestWait_ReturnsSignal(t *testing.T) {
	notified := make(chan os.Signal, 1)
	h := &Handler{
		signals: []os.Signal{syscall.SIGTERM},
		notify: func(ch chan<- os.Signal, sigs ...os.Signal) {
			go func() { ch <- syscall.SIGTERM }()
		},
		stop: func(chan<- os.Signal) {},
	}

	go func() { notified <- h.Wait(context.Background()) }()

	select {
	case sig := <-notified:
		if sig != syscall.SIGTERM {
			t.Fatalf("expected SIGTERM, got %v", sig)
		}
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after signal")
	}
}

func TestWithCancel_CancelsOnSignal(t *testing.T) {
	h := &Handler{
		signals: []os.Signal{syscall.SIGINT},
		notify: func(ch chan<- os.Signal, sigs ...os.Signal) {
			go func() { ch <- syscall.SIGINT }()
		},
		stop: func(chan<- os.Signal) {},
	}

	ctx, cancel := h.WithCancel(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		// success
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after signal")
	}
}

func TestWithCancel_CancelsOnParentDone(t *testing.T) {
	h := New()
	parent, parentCancel := context.WithCancel(context.Background())

	ctx, cancel := h.WithCancel(parent)
	defer cancel()

	parentCancel()

	select {
	case <-ctx.Done():
		// success — parent propagated
	case <-time.After(time.Second):
		t.Fatal("derived context not cancelled when parent cancelled")
	}
}
