// Package portpulse tracks the rate of port change events over time,
// emitting a "pulse" metric that reflects how active the monitored
// host has been during a rolling window.
package portpulse

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Pulse holds a single pulse measurement.
type Pulse struct {
	At      time.Time
	Opened  int
	Closed  int
	Total   int
}

// Tracker accumulates port-change events and reports a rolling pulse.
type Tracker struct {
	mu     sync.Mutex
	window time.Duration
	events []Pulse
	out    io.Writer
}

// New returns a Tracker with the given rolling window.
// If out is nil it defaults to os.Stdout.
func New(window time.Duration, out io.Writer) *Tracker {
	if out == nil {
		out = os.Stdout
	}
	return &Tracker{
		window: window,
		out:    out,
	}
}

// Observe records opened/closed counts derived from a snapshot diff and
// prunes events that have fallen outside the rolling window.
func (t *Tracker) Observe(diff snapshot.Diff) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	p := Pulse{
		At:     now,
		Opened: len(diff.Opened),
		Closed: len(diff.Closed),
		Total:  len(diff.Opened) + len(diff.Closed),
	}
	t.events = append(t.events, p)
	t.prune(now)
}

// Rate returns the total number of change events within the rolling window.
func (t *Tracker) Rate() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(time.Now())
	total := 0
	for _, e := range t.events {
		total += e.Total
	}
	return total
}

// Report writes a human-readable pulse summary to the configured writer.
func (t *Tracker) Report() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(time.Now())
	rate := 0
	for _, e := range t.events {
		rate += e.Total
	}
	fmt.Fprintf(t.out, "portpulse: %d change events in last %s\n", rate, t.window)
}

// prune removes events older than the rolling window. Caller must hold mu.
func (t *Tracker) prune(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.events) && t.events[i].At.Before(cutoff) {
		i++
	}
	t.events = t.events[i:]
}
