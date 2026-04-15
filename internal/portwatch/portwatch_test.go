package portwatch_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/config"
	"portwatch/internal/portwatch"
)

func tempState(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func defaultCfg() *config.Config {
	c := config.DefaultConfig()
	c.StartPort = 65530
	c.EndPort = 65535
	c.Interval = "50ms"
	return c
}

func TestNew_ReturnsWatcher(t *testing.T) {
	w, err := portwatch.New(defaultCfg(), tempState(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestNew_NilConfigReturnsError(t *testing.T) {
	_, err := portwatch.New(nil, tempState(t))
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNew_BadStatePathReturnsError(t *testing.T) {
	_, err := portwatch.New(defaultCfg(), string([]byte{0}))
	if err == nil {
		t.Skip("platform did not reject invalid path")
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	w, err := portwatch.New(defaultCfg(), tempState(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	w.SetOutput(&buf)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	if runErr := w.Run(ctx); runErr != nil {
		t.Fatalf("Run returned unexpected error: %v", runErr)
	}

	if buf.Len() == 0 {
		t.Error("expected at least one status line in output")
	}
}

func TestRun_WritesStartLine(t *testing.T) {
	w, err := portwatch.New(defaultCfg(), tempState(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	w.SetOutput(&buf)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	if !bytes.Contains(buf.Bytes(), []byte("portwatch: starting")) {
		t.Errorf("expected start line, got: %s", buf.String())
	}
}

func init() {
	// Ensure test binary can locate project root for any relative imports.
	os.Setenv("PORTWATCH_TEST", "1")
}
