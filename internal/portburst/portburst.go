// Package portburst detects sudden bursts of new ports opening within a
// short observation window and emits a warning when the count exceeds a
// configurable threshold.
package portburst

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Detector tracks port-open events and warns on bursts.
type Detector struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	events    []time.Time
	out        io.Writer
}

// New returns a Detector that fires when more than threshold ports are opened
// within window. Output defaults to os.Stderr when w is nil.
func New(threshold int, window time.Duration, w io.Writer) *Detector {
	if w == nil {
		w = os.Stderr
	}
	if threshold <= 0 {
		threshold = 5
	}
	if window <= 0 {
		window = 30 * time.Second
	}
	return &Detector{
		threshold: threshold,
		window:    window,
		out:        w,
	}
}

// Observe records the opened entries from diff and warns if the burst
// threshold is breached within the configured window.
func (d *Detector) Observe(diff snapshot.Diff) {
	if len(diff.Opened) == 0 {
		return
	}
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()

	// Append one timestamp per opened port.
	for range diff.Opened {
		d.events = append(d.events, now)
	}

	// Prune events outside the window.
	cutoff := now.Add(-d.window)
	valid := d.events[:0]
	for _, t := range d.events {
		if !t.Before(cutoff) {
			valid = append(valid, t)
		}
	}
	d.events = valid

	if len(d.events) >= d.threshold {
		fmt.Fprintf(d.out, "[portburst] burst detected: %d ports opened within %s\n",
			len(d.events), d.window)
	}
}

// Reset clears the internal event history.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = d.events[:0]
}

// Count returns the number of events currently within the window.
func (d *Detector) Count(now time.Time) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	cutoff := now.Add(-d.window)
	count := 0
	for _, t := range d.events {
		if !t.Before(cutoff) {
			count++
		}
	}
	return count
}
