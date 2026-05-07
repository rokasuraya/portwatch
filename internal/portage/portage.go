// Package portage tracks how long each port has been continuously open.
package portage

import (
	"fmt"
	"sync"
	"time"

	"portwatch/internal/snapshot"
)

// Entry holds age metadata for a single port.
type Entry struct {
	Port      int
	Protocol  string
	FirstSeen time.Time
	LastSeen  time.Time
	SeenCount int
}

// Age returns how long the port has been continuously observed.
func (e Entry) Age(now time.Time) time.Duration {
	return now.Sub(e.FirstSeen)
}

// Tracker records first-seen and last-seen times for open ports.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
	now     func() time.Time
}

// New returns a new Tracker. If now is nil, time.Now is used.
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{
		entries: make(map[string]*Entry),
		now:     now,
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Observe updates age tracking from the provided snapshot.
func (t *Tracker) Observe(snap *snapshot.Snapshot) {
	if snap == nil {
		return
	}
	now := t.now()
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[string]struct{})
	for _, e := range snap.Entries {
		k := portKey(e.Port, e.Protocol)
		seen[k] = struct{}{}
		if existing, ok := t.entries[k]; ok {
			existing.LastSeen = now
			existing.SeenCount++
		} else {
			t.entries[k] = &Entry{
				Port:      e.Port,
				Protocol:  e.Protocol,
				FirstSeen: now,
				LastSeen:  now,
				SeenCount: 1,
			}
		}
	}
	// Remove ports no longer present.
	for k := range t.entries {
		if _, ok := seen[k]; !ok {
			delete(t.entries, k)
		}
	}
}

// Get returns the Entry for the given port/protocol, and whether it exists.
func (t *Tracker) Get(port int, proto string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[portKey(port, proto)]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a copy of all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}
