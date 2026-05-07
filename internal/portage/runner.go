package portage

import (
	"context"
	"log"
	"time"

	"portwatch/internal/snapshot"
)

// SnapshotFunc returns the current snapshot or an error.
type SnapshotFunc func() (*snapshot.Snapshot, error)

// Runner periodically feeds snapshots into a Tracker.
type Runner struct {
	tracker  *Tracker
	snap     SnapshotFunc
	interval time.Duration
}

// NewRunner returns a Runner that calls snap on each tick.
func NewRunner(tracker *Tracker, snap SnapshotFunc, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Runner{
		tracker:  tracker,
		snap:     snap,
		interval: interval,
	}
}

// Run starts the observation loop. It blocks until ctx is cancelled.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s, err := r.snap()
			if err != nil {
				log.Printf("portage: snapshot error: %v", err)
				continue
			}
			r.tracker.Observe(s)
		}
	}
}
