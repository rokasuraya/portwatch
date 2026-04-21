package portage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

const defaultInterval = 30 * time.Second

// SnapshotFunc returns the current snapshot or an error.
type SnapshotFunc func() (*snapshot.Snapshot, error)

// Runner periodically drives the Tracker's Observe method using a
// caller-supplied snapshot function.
type Runner struct {
	tracker  *Tracker
	snap     SnapshotFunc
	interval time.Duration
	log      *log.Logger
}

// NewRunner returns a Runner that calls snap on each tick and feeds the
// result into tracker. If interval is zero, defaultInterval is used.
func NewRunner(tracker *Tracker, snap SnapshotFunc, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = defaultInterval
	}
	return &Runner{
		tracker:  tracker,
		snap:     snap,
		interval: interval,
		log:      log.New(os.Stderr, "portage/runner: ", 0),
	}
}

// Run blocks until ctx is cancelled, calling Observe on each tick.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snap, err := r.snap()
			if err != nil {
				r.log.Printf("snapshot error: %v", err)
				continue
			}
			r.tracker.Observe(snap)
		}
	}
}
