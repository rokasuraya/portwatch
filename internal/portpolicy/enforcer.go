package portpolicy

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"portwatch/internal/snapshot"
)

// SnapshotFunc returns the current snapshot or an error.
type SnapshotFunc func() (*snapshot.Snapshot, error)

// Enforcer periodically evaluates a Policy against live snapshots and
// writes violations to an output writer.
type Enforcer struct {
	policy   *Policy
	snap     SnapshotFunc
	interval time.Duration
	out      io.Writer
}

// NewEnforcer creates an Enforcer. If out is nil, os.Stderr is used.
func NewEnforcer(p *Policy, fn SnapshotFunc, interval time.Duration, out io.Writer) *Enforcer {
	if out == nil {
		out = os.Stderr
	}
	return &Enforcer{
		policy:   p,
		snap:     fn,
		interval: interval,
		out:      out,
	}
}

// Run starts the enforcement loop. It blocks until ctx is cancelled.
func (e *Enforcer) Run(ctx context.Context) error {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			e.enforce()
		}
	}
}

func (e *Enforcer) enforce() {
	snap, err := e.snap()
	if err != nil {
		fmt.Fprintf(e.out, "portpolicy: snapshot error: %v\n", err)
		return
	}
	for _, v := range e.policy.Check(snap) {
		fmt.Fprintf(e.out, "portpolicy violation: %s\n", v.Error())
	}
}
