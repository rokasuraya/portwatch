package pipeline_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/snapshot"
)

// TestIntegration_MultipleTicksAccumulateDiff verifies that running several
// ticks with a changing scan result correctly accumulates opened/closed diffs
// across stages.
func TestIntegration_MultipleTicksAccumulateDiff(t *testing.T) {
	scans := [][]snapshot.Entry{
		{{Port: 22, Proto: "tcp"}},
		{{Port: 22, Proto: "tcp"}, {Port: 80, Proto: "tcp"}},
		{{Port: 80, Proto: "tcp"}},
	}

	call := 0
	scan := func(_ context.Context) ([]snapshot.Entry, error) {
		defer func() { call++ }()
		if call >= len(scans) {
			return scans[len(scans)-1], nil
		}
		return scans[call], nil
	}

	var totalOpened, totalClosed int64
	stage := func(_ context.Context, opened, closed []snapshot.Entry) error {
		atomic.AddInt64(&totalOpened, int64(len(opened)))
		atomic.AddInt64(&totalClosed, int64(len(closed)))
		return nil
	}

	p := pipeline.New(scan, makeStore(t), stage)
	for i := 0; i < len(scans); i++ {
		if _, err := p.Tick(context.Background()); err != nil {
			t.Fatalf("tick %d: %v", i, err)
		}
	}

	// tick 0: no previous → no diff
	// tick 1: port 80 opened
	// tick 2: port 22 closed
	if totalOpened != 1 {
		t.Errorf("expected 1 opened event, got %d", totalOpened)
	}
	if totalClosed != 1 {
		t.Errorf("expected 1 closed event, got %d", totalClosed)
	}
}

// TestIntegration_ContextCancelledDuringScan ensures a cancelled context
// propagates correctly through the pipeline.
func TestIntegration_ContextCancelledDuringScan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	scan := func(c context.Context) ([]snapshot.Entry, error) {
		if err := c.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	p := pipeline.New(scan, makeStore(t))
	_, err := p.Tick(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}
