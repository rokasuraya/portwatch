// Package window provides a sliding time-window counter for tracking
// event frequencies over a rolling duration.
package window

import (
	"sync"
	"time"
)

// Window tracks how many events have occurred within a rolling time window.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	buckets  []bucket
	size     int
}

type bucket struct {
	at    time.Time
	count int
}

// New creates a Window that retains events within the given duration.
// size controls the number of time buckets used for granularity.
func New(duration time.Duration, size int) *Window {
	if size < 1 {
		size = 1
	}
	return &Window{
		duration: duration,
		size:     size,
		buckets:  make([]bucket, 0, size),
	}
}

// Add records n events at the current time.
func (w *Window) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	w.evict(now)
	w.buckets = append(w.buckets, bucket{at: now, count: n})
}

// Count returns the total number of events within the current window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	total := 0
	for _, b := range w.buckets {
		total += b.count
	}
	return total
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = w.buckets[:0]
}

// Oldest returns the timestamp of the earliest event still within the window,
// and false if the window is empty.
func (w *Window) Oldest() (time.Time, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	if len(w.buckets) == 0 {
		return time.Time{}, false
	}
	return w.buckets[0].at, true
}

// evict removes buckets that have fallen outside the window. Must be called
// with w.mu held.
func (w *Window) evict(now time.Time) {
	cutoff := now.Add(-w.duration)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
