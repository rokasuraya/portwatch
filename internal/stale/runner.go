package stale

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// SnapshotFunc is a function that returns the current snapshot.
type SnapshotFunc func(ctx context.Context) (*snapshot.Snapshot, error)

// Runner periodically calls a SnapshotFunc, feeds results into a Detector,
// and then calls Check to emit stale-port warnings.
type Runner struct {
	detector *Detector
	getSnap  SnapshotFunc
	interval time.Duration
}

// NewRunner creates a Runner that polls every interval.
func NewRunner(d *Detector, fn SnapshotFunc, interval time.Duration) *Runner {
	return &Runner{detector: d, getSnap: fn, interval: interval}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.tick(ctx)
		}
	}
}

// tick performs a single observe+check cycle.
func (r *Runner) tick(ctx context.Context) {
	snap, err := r.getSnap(ctx)
	if err != nil {
		return
	}
	r.detector.Observe(snap)
	r.detector.Check()
}
