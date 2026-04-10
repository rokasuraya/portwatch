// Package throttle provides alert rate-limiting to prevent notification floods
// when a port repeatedly opens and closes within a short window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks recent alert events per key and suppresses duplicates
// that occur within the configured cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New returns a Throttle with the given cooldown duration.
// Alerts for the same key are suppressed until the cooldown has elapsed.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether an alert for the given key should be forwarded.
// It returns true the first time a key is seen, and again only after the
// cooldown window has elapsed since the last allowed alert.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.last[key]; ok && now.Sub(last) < t.cooldown {
		return false
	}
	t.last[key] = now
	return true
}

// Reset clears the throttle state for the given key, allowing the next
// alert for that key to pass through immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Purge removes all entries whose last-seen time is older than the cooldown,
// keeping memory usage bounded during long-running daemon sessions.
func (t *Throttle) Purge() {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := time.Now().Add(-t.cooldown)
	for k, ts := range t.last {
		if ts.Before(cutoff) {
			delete(t.last, k)
		}
	}
}
