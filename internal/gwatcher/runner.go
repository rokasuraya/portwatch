package gwatcher

import (
	"context"
	"time"

	"portwatch/internal/snapshot"
)

// SnapshotFunc returns the current snapshot or an error.
type SnapshotFunc func() (*snapshot.Snapshot, error)

// OnEventFunc is called with each batch of group-change events.
type OnEventFunc func([]Event)

// Runner periodically calls a SnapshotFunc and forwards group-change events.
type Runner struct {
	watcher  *Watcher
	snap     SnapshotFunc
	onEvent  OnEventFunc
	interval time.Duration
}

// NewRunner constructs a Runner.
func NewRunner(w *Watcher, snap SnapshotFunc, onEvent OnEventFunc, interval time.Duration) *Runner {
	return &Runner{watcher: w, snap: snap, onEvent: onEvent, interval: interval}
}

// Run blocks until ctx is cancelled, polling at the configured interval.
func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s, err := r.snap()
			if err != nil {
				continue
			}
			if events := r.watcher.Observe(s); len(events) > 0 {
				r.onEvent(events)
			}
		}
	}
}
