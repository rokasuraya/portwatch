package supervisor_test

import (
	"context"
	"testing"
	"time"

	"portwatch/internal/config"
	"portwatch/internal/metrics"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
	"portwatch/internal/supervisor"
)

// TestIntegration_MetricsUpdatedAfterTick verifies that a completed tick
// increments the scan counter inside the Metrics component.
func TestIntegration_MetricsUpdatedAfterTick(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.PortRangeStart = 1
	cfg.PortRangeEnd = 5

	sc := scanner.New(50 * time.Millisecond)
	st, err := state.New(tempStateFile(t))
	if err != nil {
		t.Fatal(err)
	}
	m := metrics.New()

	sv := supervisor.New(cfg, supervisor.Components{
		Scanner: sc,
		State:   st,
		Metrics: m,
	}, nil)

	if err := sv.Tick(context.Background()); err != nil {
		t.Fatalf("tick error: %v", err)
	}

	snap := m.Snapshot()
	if snap.ScanCount != 1 {
		t.Errorf("expected ScanCount=1, got %d", snap.ScanCount)
	}
	if snap.LastScanDuration <= 0 {
		t.Error("expected positive LastScanDuration")
	}
}

// TestIntegration_TwoTicksDetectsStableState verifies that running two
// consecutive ticks on an unchanging host produces zero diffs on the
// second tick (opened and closed remain empty).
func TestIntegration_TwoTicksDetectsStableState(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.PortRangeStart = 1
	cfg.PortRangeEnd = 5

	sc := scanner.New(50 * time.Millisecond)
	st, err := state.New(tempStateFile(t))
	if err != nil {
		t.Fatal(err)
	}
	m := metrics.New()

	changes := 0
	sv := supervisor.New(cfg, supervisor.Components{
		Scanner: sc,
		State:   st,
		Metrics: m,
		OnChange: func(_, _ []string) { changes++ },
	}, nil)

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		if err := sv.Tick(ctx); err != nil {
			t.Fatalf("tick %d error: %v", i, err)
		}
	}

	// The second tick should not trigger onChange because state is stable.
	if changes > 1 {
		t.Errorf("expected at most 1 onChange call, got %d", changes)
	}
}
