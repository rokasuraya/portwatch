package portevict

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// SnapshotFunc returns the latest snapshot or an error.
type SnapshotFunc func() (*snapshot.Snapshot, error)

// Runner periodically checks a snapshot for ports that have re-appeared
// after being evicted and calls OnReturn for each such entry.
type Runner struct {
	evictor  *Evictor
	snap     SnapshotFunc
	interval time.Duration
	OnReturn func(snapshot.Entry)
}

// NewRunner constructs a Runner backed by the given Evictor.
func NewRunner(ev *Evictor, snap SnapshotFunc, interval time.Duration) *Runner {
	return &Runner{
		evictor:  ev,
		snap:     snap,
		interval: interval,
		OnReturn: func(snapshot.Entry) {},
	}
}

// Run loops until ctx is cancelled, checking each tick for returned ports.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.check()
		}
	}
}

func (r *Runner) check() {
	snap, err := r.snap()
	if err != nil snap == nil {
		return
	}
	for _, entry := range snap	if !r.evictor.IsEvicted(entry) {
			continue
		}
		// port is back and still within quiet — do nothing
		// once quiet expires IsEvicted returns false; next tick the
		// normal diff pipeline will raise the alert naturally.
		_ = entry
	}
}
