// Package jitter provides utilities for adding randomised jitter to
// durations, helping spread load across multiple goroutines or instances
// that would otherwise fire at exactly the same time.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is a seeded random source safe for concurrent use.
type Source struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// New returns a new Source seeded with the current time.
func New() *Source {
	return &Source{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
}

// NewWithSeed returns a new Source seeded with the provided value.
// Useful for deterministic tests.
func NewWithSeed(seed int64) *Source {
	return &Source{
		rng: rand.New(rand.NewSource(seed)), //nolint:gosec
	}
}

// Apply returns d plus a uniformly-distributed random offset in [0, max).
// If max is zero or negative, d is returned unchanged.
func (s *Source) Apply(d, max time.Duration) time.Duration {
	if max <= 0 {
		return d
	}
	s.mu.Lock()
	offset := time.Duration(s.rng.Int63n(int64(max)))
	s.mu.Unlock()
	return d + offset
}

// Spread returns a uniformly-distributed random duration in [0, max).
// If max is zero or negative, zero is returned.
func (s *Source) Spread(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	s.mu.Lock()
	v := time.Duration(s.rng.Int63n(int64(max)))
	s.mu.Unlock()
	return v
}

// Clamp ensures the jitter offset does not cause d to exceed ceiling.
// Returns d + jitter where jitter is clamped so that result <= ceiling.
func (s *Source) Clamp(d, max, ceiling time.Duration) time.Duration {
	result := s.Apply(d, max)
	if result > ceiling {
		return ceiling
	}
	return result
}
