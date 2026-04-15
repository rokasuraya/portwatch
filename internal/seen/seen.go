// Package seen tracks which ports have been observed across scans,
// providing a simple first-seen / last-seen ledger keyed by protocol and port.
package seen

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry records the first and most-recent observation of a port.
type Entry struct {
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

// Ledger maintains a thread-safe map of port observations.
type Ledger struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	now     func() time.Time
}

// New returns an initialised Ledger.
func New() *Ledger {
	return &Ledger{
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

func portKey(e snapshot.Entry) string {
	return fmt.Sprintf("%s:%d", e.Protocol, e.Port)
}

// Observe records every entry in the snapshot as seen at the current time.
func (l *Ledger) Observe(snap *snapshot.Snapshot) {
	if snap == nil {
		return
	}
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range snap.Entries {
		k := portKey(e)
		if existing, ok := l.entries[k]; ok {
			existing.LastSeen = now
			existing.Count++
		} else {
			l.entries[k] = &Entry{
				FirstSeen: now,
				LastSeen:  now,
				Count:     1,
			}
		}
	}
}

// Lookup returns the Entry for the given snapshot entry, and whether it exists.
func (l *Ledger) Lookup(e snapshot.Entry) (Entry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	v, ok := l.entries[portKey(e)]
	if !ok {
		return Entry{}, false
	}
	return *v, true
}

// Len returns the number of distinct ports tracked.
func (l *Ledger) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}

// Reset clears all tracked entries.
func (l *Ledger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make(map[string]*Entry)
}
