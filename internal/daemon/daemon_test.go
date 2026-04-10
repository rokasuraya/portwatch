package daemon

import (
	"context"
	"os"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/state"
)

func defaultConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	// Use a very small port range to keep tests fast.
	cfg.PortStart = 1
	cfg.PortEnd = 10
	cfg.Interval = "100ms"
	return cfg
}

func tempStateFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-state-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()
	os.Remove(f.Name()) // state.New creates it fresh
	return f.Name()
}

func TestNew_ReturnsDaemon(t *testing.T) {
	cfg := defaultConfig(t)
	st, err := state.New(tempStateFile(t))
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}
	al := alert.New(nil)
	d := New(cfg, st, al)
	if d == nil {
		t.Fatal("expected non-nil daemon")
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	cfg := defaultConfig(t)
	st, err := state.New(tempStateFile(t))
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}
	al := alert.New(nil)
	d := New(cfg, st, al)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestTick_DoesNotError(t *testing.T) {
	cfg := defaultConfig(t)
	st, err := state.New(tempStateFile(t))
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}
	al := alert.New(nil)
	d := New(cfg, st, al)

	if err := d.tick(); err != nil {
		t.Fatalf("tick returned unexpected error: %v", err)
	}
}
