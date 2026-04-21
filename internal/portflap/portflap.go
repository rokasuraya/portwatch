// Package portflap detects ports that open and close rapidly (flapping)
// within a configurable observation window, helping distinguish transient
// noise from stable state changes.
package portflap

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Detector tracks open/close transitions per port and warns when a port
// flaps more than Threshold times within Window.
type Detector struct {
	mu        sync.Mutex
	counts    map[string][]time.Time
	Threshold int
	Window    time.Duration
	out       io.Writer
	now       func() time.Time
}

// New returns a Detector with the given threshold and window.
// Warnings are written to stderr by default.
func New(threshold int, window time.Duration) *Detector {
	return &Detector{
		counts:    make(map[string][]time.Time),
		Threshold: threshold,
		Window:    window,
		out:       os.Stderr,
		now:       time.Now,
	}
}

// SetOutput redirects warning output.
func (d *Detector) SetOutput(w io.Writer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.out = w
}

func portKey(e snapshot.Entry) string {
	return fmt.Sprintf("%d/%s", e.Port, e.Protocol)
}

// Observe records transitions from the diff and emits a warning for any
// port whose flap count exceeds the threshold within the window.
func (d *Detector) Observe(opened, closed []snapshot.Entry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	cutoff := now.Add(-d.Window)

	record := func(e snapshot.Entry) {
		k := portKey(e)
		times := d.counts[k]
		// prune old events
		filtered := times[:0]
		for _, t := range times {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		filtered = append(filtered, now)
		d.counts[k] = filtered
		if len(filtered) >= d.Threshold {
			fmt.Fprintf(d.out, "portflap: port %s flapped %d times in %s\n",
				k, len(filtered), d.Window)
		}
	}

	for _, e := range opened {
		record(e)
	}
	for _, e := range closed {
		record(e)
	}
}

// Reset clears all recorded transitions.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.counts = make(map[string][]time.Time)
}
