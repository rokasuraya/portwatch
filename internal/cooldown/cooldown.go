// Package cooldown provides a per-key cooldown tracker that suppresses
// repeated events within a configurable quiet period.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last-seen time for arbitrary string keys and reports
// whether enough time has elapsed to allow a new event through.
type Cooldown struct {
	mu      sync.Mutex
	period  time.Duration
	lastSeen map[string]time.Time
	now     func() time.Time
}

// New returns a Cooldown with the given quiet period.
func New(period time.Duration) *Cooldown {
	return &Cooldown{
		period:   period,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the key has not been seen within the cooldown period
// and records the current time as the last-seen time for that key.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if last, ok := c.lastSeen[key]; ok {
		if now.Sub(last) < c.period {
			return false
		}
	}
	c.lastSeen[key] = now
	return true
}

// Reset removes the cooldown record for the given key, allowing the next
// call to Allow to pass unconditionally.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastSeen, key)
}

// Len returns the number of keys currently tracked.
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.lastSeen)
}
