// Package dedupe provides a lightweight deduplication filter that suppresses
// repeated scanner entries within a configurable time window.
package dedupe

import (
	"fmt"
	"sync"
	"time"
)

// Entry represents a single port observation that may be deduplicated.
type Entry struct {
	Host     string
	Port     uint16
	Protocol string
}

func key(e Entry) string {
	return fmt.Sprintf("%s/%s/%d", e.Host, e.Protocol, e.Port)
}

// Filter suppresses duplicate entries seen within Window.
type Filter struct {
	mu     sync.Mutex
	seen   map[string]time.Time
	window time.Duration
	now    func() time.Time
}

// New returns a Filter that deduplicates entries within the given window.
// A zero or negative window disables deduplication (every entry is allowed).
func New(window time.Duration) *Filter {
	return &Filter{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// Allow returns true if the entry has not been seen within the current window.
// If the entry is new (or the window has expired), it is recorded and true is
// returned. Otherwise false is returned and the entry is suppressed.
func (f *Filter) Allow(e Entry) bool {
	if f.window <= 0 {
		return true
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.now()
	k := key(e)
	if t, ok := f.seen[k]; ok && now.Sub(t) < f.window {
		return false
	}
	f.seen[k] = now
	return true
}

// Flush removes all recorded entries, resetting the filter state.
func (f *Filter) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seen = make(map[string]time.Time)
}

// Len returns the number of entries currently tracked.
func (f *Filter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}
