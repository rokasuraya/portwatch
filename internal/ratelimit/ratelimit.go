// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently port scan cycles may be triggered externally.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a simple token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	rate     time.Duration // minimum interval between allowed events
	last     time.Time
	bursts   int
	tokens   int
	maxBurst int
}

// New creates a Limiter that allows at most one event per interval,
// with an optional burst capacity (minimum 1).
func New(interval time.Duration, burst int) *Limiter {
	if burst < 1 {
		burst = 1
	}
	return &Limiter{
		rate:     interval,
		tokens:   burst,
		maxBurst: burst,
	}
}

// Allow reports whether an event is permitted right now.
// It consumes one token if available, refilling based on elapsed time.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if !l.last.IsZero() {
		elapsed := now.Sub(l.last)
		refill := int(elapsed / l.rate)
		if refill > 0 {
			l.tokens += refill
			if l.tokens > l.maxBurst {
				l.tokens = l.maxBurst
			}
			l.last = l.last.Add(time.Duration(refill) * l.rate)
		}
	} else {
		l.last = now
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// Reset restores the limiter to its full burst capacity.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.maxBurst
	l.last = time.Time{}
}

// Remaining returns the number of tokens currently available.
func (l *Limiter) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
