// Package portcap tracks the capacity of observed ports over time,
// recording the maximum number of simultaneously open ports seen per
// protocol and emitting a warning when a new high-water mark is reached.
package portcap

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Record holds the high-water mark for a single protocol.
type Record struct {
	Protocol string
	Peak     int
}

// Tracker monitors open-port counts and reports new peak values.
type Tracker struct {
	mu     sync.Mutex
	peaks  map[string]int // keyed by protocol
	out    io.Writer
}

// New returns a Tracker that writes peak-breach notices to out.
// If out is nil, os.Stderr is used.
func New(out io.Writer) *Tracker {
	if out == nil {
		out = os.Stderr
	}
	return &Tracker{
		peaks: make(map[string]int),
		out:   out,
	}
}

// Observe inspects snap and updates the high-water mark for each
// protocol present. When a new peak is reached a notice is written to
// the configured writer. Observe is safe for concurrent use.
func (t *Tracker) Observe(snap *snapshot.Snapshot) {
	if snap == nil {
		return
	}

	// Count open ports per protocol.
	counts := make(map[string]int)
	for _, e := range snap.Entries {
		counts[e.Protocol]++
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for proto, count := range counts {
		if prev, ok := t.peaks[proto]; !ok || count > prev {
			t.peaks[proto] = count
			fmt.Fprintf(t.out, "portcap: new peak for %s — %d open ports\n", proto, count)
		}
	}
}

// Peak returns the recorded high-water mark for the given protocol.
// It returns 0 if no observation has been made for that protocol.
func (t *Tracker) Peak(protocol string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.peaks[protocol]
}

// Records returns a snapshot of all recorded peaks, one entry per
// protocol. The slice is sorted by protocol name for determinism.
func (t *Tracker) Records() []Record {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make([]Record, 0, len(t.peaks))
	for proto, peak := range t.peaks {
		out = append(out, Record{Protocol: proto, Peak: peak})
	}
	// Simple insertion sort — the number of protocols is tiny.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Protocol < out[j-1].Protocol; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// Reset clears all recorded peaks. Useful in tests or when the
// monitoring target changes completely.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peaks = make(map[string]int)
}
