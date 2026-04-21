package portcap

import (
	"context"
	"fmt"
	"time"
)

// SnapshotFunc returns the current snapshot or an error.
type SnapshotFunc func() (*Snapshot, error)

// Runner periodically checks the port cap and writes violations.
type Runner struct {
	cap      *PortCap
	snapshotFn SnapshotFunc
	interval time.Duration
}

// NewRunner creates a Runner that checks the cap on every interval.
func NewRunner(c *PortCap, fn SnapshotFunc, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Runner{cap: c, snapshotFn: fn, interval: interval}
}

// Run starts the periodic cap check loop. It blocks until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			snap, err := r.snapshotFn()
			if err != nil {
				fmt.Fprintf(r.cap.w, "portcap: snapshot error: %v\n", err)
				continue
			}
			r.cap.Check(snap)
		}
	}
}
