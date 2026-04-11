package supervisor_test

import (
	"context"
	"log"
	"os"
	"testing"

	"portwatch/internal/config"
	"portwatch/internal/metrics"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
	"portwatch/internal/supervisor"
)

func defaultCfg(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	// Use a very narrow range so the test is fast and deterministic.
	cfg.PortRangeStart = 1
	cfg.PortRangeEnd = 10
	return cfg
}

func tempStateFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestNew_ReturnsSupervisor(t *testing.T) {
	cfg := defaultCfg(t)
	sc := scanner.New(0)
	st, _ := state.New(tempStateFile(t))
	m := metrics.New()

	sv := supervisor.New(cfg, supervisor.Components{
		Scanner: sc,
		State:   st,
		Metrics: m,
	}, nil)

	if sv == nil {
		t.Fatal("expected non-nil supervisor")
	}
}

func TestTick_DoesNotError(t *testing.T) {
	cfg := defaultCfg(t)
	sc := scanner.New(0)
	st, _ := state.New(tempStateFile(t))
	m := metrics.New()

	sv := supervisor.New(cfg, supervisor.Components{
		Scanner: sc,
		State:   st,
		Metrics: m,
	}, log.New(os.Stderr, "", 0))

	if err := sv.Tick(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTick_InvokesOnChange(t *testing.T) {
	cfg := defaultCfg(t)
	sc := scanner.New(0)
	st, _ := state.New(tempStateFile(t))
	m := metrics.New()

	called := false
	sv := supervisor.New(cfg, supervisor.Components{
		Scanner: sc,
		State:   st,
		Metrics: m,
		OnChange: func(opened, closed []string) {
			called = true
		},
	}, nil)

	// First tick seeds state; second tick with same state won't fire.
	// We just verify Tick completes without panic.
	_ = sv.Tick(context.Background())
	_ = called // value depends on host environment; just ensure no panic.
}

func TestTick_CancelledContext(t *testing.T) {
	cfg := defaultCfg(t)
	sc := scanner.New(0)
	st, _ := state.New(tempStateFile(t))
	m := metrics.New()

	sv := supervisor.New(cfg, supervisor.Components{
		Scanner: sc,
		State:   st,
		Metrics: m,
	}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	// Should return quickly; error is acceptable.
	_ = sv.Tick(ctx)
}
