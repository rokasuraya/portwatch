// Package watchdog provides a self-monitoring component that detects when
// the scan loop has stalled or missed expected ticks, and emits alerts.
package watchdog

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Watchdog monitors scan heartbeats and warns when ticks are missed.
type Watchdog struct {
	mu        sync.Mutex
	lastBeat  time.Time
	tolerance time.Duration
	out        io.Writer
}

// New returns a Watchdog that emits a warning to out (defaults to os.Stderr)
// if no heartbeat is received within tolerance of the expected interval.
func New(tolerance time.Duration, out io.Writer) *Watchdog {
	if out == nil {
		out = os.Stderr
	}
	return &Watchdog{
		tolerance: tolerance,
		out:       out,
	}
}

// Beat records a heartbeat at the current time.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastBeat = time.Now()
}

// Run starts the watchdog loop, checking every interval whether a beat has
// been received within tolerance. It blocks until ctx is cancelled.
func (w *Watchdog) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Seed the first beat so we don't false-alarm on startup.
	w.Beat()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			w.check(t, interval)
		}
	}
}

// check compares the last beat time against the deadline and writes a warning
// when the deadline is exceeded.
func (w *Watchdog) check(now time.Time, interval time.Duration) {
	w.mu.Lock()
	last := w.lastBeat
	w.mu.Unlock()

	deadline := last.Add(interval + w.tolerance)
	if now.After(deadline) {
		since := now.Sub(last).Round(time.Millisecond)
		fmt.Fprintf(w.out, "[watchdog] WARNING: no scan beat for %s (tolerance %s)\n",
			since, w.tolerance)
	}
}
