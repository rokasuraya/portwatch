// Package limiter provides a concurrency limiter that caps the number of
// simultaneous port-scan goroutines, preventing resource exhaustion on large
// port ranges.
package limiter

import "context"

// Limiter controls the maximum number of concurrent workers.
type Limiter struct {
	sem chan struct{}
}

// New returns a Limiter that allows at most n concurrent acquisitions.
// If n <= 0 it defaults to 1.
func New(n int) *Limiter {
	if n <= 0 {
		n = 1
	}
	return &Limiter{sem: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns ctx.Err() if the context is done before a slot is obtained.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot.
// It panics if called more times than Acquire has succeeded.
func (l *Limiter) Release() {
	select {
	case <-l.sem:
	default:
		panic("limiter: Release called without a matching Acquire")
	}
}

// Cap returns the maximum concurrency configured for this Limiter.
func (l *Limiter) Cap() int {
	return cap(l.sem)
}

// Available returns the number of slots that can be acquired right now
// without blocking.
func (l *Limiter) Available() int {
	return cap(l.sem) - len(l.sem)
}
