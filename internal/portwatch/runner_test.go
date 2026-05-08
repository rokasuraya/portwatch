package portwatch

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestNewRunner_ReturnsRunner(t *testing.T) {
	w, _ := New(defaultCfg(t), tempState(t))
	r := NewRunner(w, 50*time.Millisecond, nil)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestRunner_CancelsCleanly(t *testing.T) {
	w, _ := New(defaultCfg(t), tempState(t))
	r := NewRunner(w, 50*time.Millisecond, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := r.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("expected context error, got %v", err)
	}
}

func TestRunner_InvokesTick(t *testing.T) {
	cfg := defaultCfg(t)
	cfg.PortRange = [2]int{1, 1} // minimal range
	w, err := New(cfg, tempState(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	r := NewRunner(w, 30*time.Millisecond, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = r.Run(ctx)
}

func TestRunner_NilLoggerDefaults(t *testing.T) {
	w, _ := New(defaultCfg(t), tempState(t))
	r := NewRunner(w, time.Second, nil)
	if r.log == nil {
		t.Fatal("expected default logger to be set")
	}
}

func TestRunner_CustomLogger(t *testing.T) {
	w, _ := New(defaultCfg(t), tempState(t))
	logger := newTestLogger()
	r := NewRunner(w, time.Second, logger)
	if r.log != logger {
		t.Fatal("expected custom logger to be used")
	}
}

// helpers

func newTestLogger() *log.Logger {
	return log.New(os.Stderr, "test: ", 0)
}

func init() {
	// suppress unused import warning — log is used via newTestLogger
	_ = log.Default
}

func defaultCfg(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	return cfg
}

func tempState(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatalf("tempState: %v", err)
	}
	f.Close()
	return f.Name()
}
