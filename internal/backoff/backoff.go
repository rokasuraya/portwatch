// Package backoff provides an exponential back-off strategy for retrying
// failed operations such as webhook deliveries or scan attempts.
package backoff

import (
	"math"
	"sync"
	"time"
)

// Backoff holds the state for a single exponential back-off sequence.
type Backoff struct {
	mu       sync.Mutex
	attempts int
	base     time.Duration
	max      time.Duration
	factor   float64
}

// New returns a Backoff with the given base delay, maximum delay, and
// multiplicative factor. Sensible defaults: base=100ms, max=30s, factor=2.0.
func New(base, max time.Duration, factor float64) *Backoff {
	if factor <= 1 {
		factor = 2.0
	}
	if base <= 0 {
		base = 100 * time.Millisecond
	}
	if max <= 0 {
		max = 30 * time.Second
	}
	return &Backoff{base: base, max: max, factor: factor}
}

// Next returns the duration to wait before the next retry and increments the
// internal attempt counter.
func (b *Backoff) Next() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	d := time.Duration(float64(b.base) * math.Pow(b.factor, float64(b.attempts)))
	if d > b.max {
		d = b.max
	}
	b.attempts++
	return d
}

// Attempts returns the number of times Next has been called.
func (b *Backoff) Attempts() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempts
}

// Reset clears the attempt counter so the sequence starts over.
func (b *Backoff) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.attempts = 0
}
