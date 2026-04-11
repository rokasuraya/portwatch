package sighandler_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sighandler"
)

func TestIntegration_Selfal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration	}

	h := sighandler.New(syscall.SIGUSR1)
	ctx := context.Background()
\tdone := make(chan os.Signal, 1)
	go func() { done <- h.Wait(ctx) }()

	// Give the goroutine time to register.
	time.Sleep(20 * time.Millisecond)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("could not find self process: %v", err)
	}
	if err := proc.Signal(syscall.SIGUSR1); err != nil {
		t.Fatalf("could not send SIGUSR1: %v", err)
	}

	select {
	case sig := <-done:
		if sig != syscall.SIGUSR1 {
			t.Fatalf("expected SIGUSR1, got %v", sig)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after self-signal")
	}
}

func TestIntegration_WithCancel_GracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := sighandler.New(syscall.SIGUSR2)
	ctx, cancel := h.WithCancel(context.Background())
	defer cancel()

	time.Sleep(20 * time.Millisecond)

	proc, _ := os.FindProcess(os.Getpid())
	_ = proc.Signal(syscall.SIGUSR2)

	select {
	case <-ctx.Done():
		// graceful
	case <-time.After(2 * time.Second):
		t.Fatal("context not cancelled after SIGUSR2")
	}
}
