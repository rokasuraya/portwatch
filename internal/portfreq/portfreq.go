// Package portfreq tracks how frequently each port appears across scans,
// providing a hit-count ledger that can be used to surface noisy or
// persistent ports over time.
package portfreq

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry holds frequency data for a single port+protocol pair.
type Entry struct {
	Port     int
	Protocol string
	Count    int64
}

// Tracker maintains per-port scan frequency counts.
type Tracker struct {
	mu      sync.Mutex
	counts  map[string]Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		counts: make(map[string]Entry),
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Observe records all open ports present in snap, incrementing each
// port's hit counter by one.
func (t *Tracker) Observe(snap *snapshot.Snapshot) {
	if snap == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, e := range snap.Entries() {
		k := portKey(e.Port, e.Protocol)
		ent := t.counts[k]
		ent.Port = e.Port
		ent.Protocol = e.Protocol
		ent.Count++
		t.counts[k] = ent
	}
}

// Get returns the frequency entry for the given port and protocol.
// ok is false when the port has never been observed.
func (t *Tracker) Get(port int, proto string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.counts[portKey(port, proto)]
	return e, ok
}

// Top returns the n ports with the highest observation counts.
// If n <= 0 all entries are returned.
func (t *Tracker) Top(n int) []Entry {
	t.mu.Lock()
	out := make([]Entry, 0, len(t.counts))
	for _, e := range t.counts {
		out = append(out, e)
	}
	t.mu.Unlock()

	// simple insertion sort — entry counts are small in practice
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Count > out[j-1].Count; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	if n > 0 && n < len(out) {
		return out[:n]
	}
	return out
}

// Reset clears all frequency data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts = make(map[string]Entry)
}
