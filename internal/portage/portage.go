// Package portage tracks how long each port has been continuously open.
package portage

import (
	"fmt"
	"sync"
	"time"

	"github.com/username/portwatch/internal/snapshot"
)

// Entry holds first-seen metadata for a single port.
type Entry struct {
	FirstSeen time.Time
	Age       time.Duration
}

// Tracker records when each port was first observed open and computes its age.
type Tracker struct {
	mu      sync.Mutex
	firstSeen map[string]time.Time
	now     func() time.Time
}

// New returns a Tracker using the real wall clock.
func New() *Tracker {
	return &Tracker{
		firstSeen: make(map[string]time.Time),
		now:       time.Now,
	}
}

func portKey(port uint16, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Observe updates the tracker from a snapshot, recording first-seen times
// for newly opened ports and removing entries for ports no longer present.
func (t *Tracker) Observe(snap *snapshot.Snapshot) {
	if snap == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	active := make(map[string]struct{}, len(snap.Entries))
	for _, e := range snap.Entries {
		k := portKey(e.Port, e.Proto)
		active[k] = struct{}{}
		if _, ok := t.firstSeen[k]; !ok {
			t.firstSeen[k] = t.now()
		}
	}

	for k := range t.firstSeen {
		if _, ok := active[k]; !ok {
			delete(t.firstSeen, k)
		}
	}
}

// Age returns the Entry for a given port/proto pair.
// ok is false if the port is not currently tracked.
func (t *Tracker) Age(port uint16, proto string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	k := portKey(port, proto)
	fs, ok := t.firstSeen[k]
	if !ok {
		return Entry{}, false
	}
	now := t.now()
	return Entry{FirstSeen: fs, Age: now.Sub(fs)}, true
}

// Len returns the number of ports currently being tracked.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.firstSeen)
}
