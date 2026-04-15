// Package trend tracks directional changes in open port counts over time,
// providing a simple rising/falling/stable signal for alerting and reporting.
package trend

import (
	"sync"
	"time"
)

// Direction represents the trend direction.
type Direction int

const (
	Stable  Direction = 0
	Rising  Direction = 1
	Falling Direction = -1
)

// String returns a human-readable label for the direction.
func (d Direction) String() string {
	switch d {
	case Rising:
		return "rising"
	case Falling:
		return "falling"
	default:
		return "stable"
	}
}

// Sample holds a single observation.
type Sample struct {
	At    time.Time
	Count int
}

// Tracker accumulates port-count samples and derives a trend direction.
type Tracker struct {
	mu      sync.Mutex
	samples []Sample
	maxLen  int
}

// New returns a Tracker that retains up to maxLen samples.
// If maxLen < 2 it is clamped to 2.
func New(maxLen int) *Tracker {
	if maxLen < 2 {
		maxLen = 2
	}
	return &Tracker{maxLen: maxLen}
}

// Record adds a new observation.
func (t *Tracker) Record(count int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples = append(t.samples, Sample{At: time.Now(), Count: count})
	if len(t.samples) > t.maxLen {
		t.samples = t.samples[len(t.samples)-t.maxLen:]
	}
}

// Direction returns the current trend by comparing the oldest and newest
// retained samples. Returns Stable when fewer than two samples exist.
func (t *Tracker) Direction() Direction {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.samples) < 2 {
		return Stable
	}
	first := t.samples[0].Count
	last := t.samples[len(t.samples)-1].Count
	switch {
	case last > first:
		return Rising
	case last < first:
		return Falling
	default:
		return Stable
	}
}

// Samples returns a copy of the current sample window.
func (t *Tracker) Samples() []Sample {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Sample, len(t.samples))
	copy(out, t.samples)
	return out
}

// Reset clears all retained samples.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.samples = t.samples[:0]
}
