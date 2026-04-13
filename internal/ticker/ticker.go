// Package ticker provides a configurable interval ticker that supports
// jitter to avoid thundering-herd problems when multiple daemons run
// concurrently.
package ticker

import (
	"context"
	"math/rand"
	"time"
)

// Ticker fires a callback on a regular interval with optional jitter.
type Ticker struct {
	interval time.Duration
	jitter   time.Duration
	onTick   func(ctx context.Context)
}

// New creates a Ticker that calls onTick every interval ± jitter.
// jitter must be less than interval; if it is not, it is clamped to zero.
func New(interval, jitter time.Duration, onTick func(ctx context.Context)) *Ticker {
	if jitter >= interval {
		jitter = 0
	}
	return &Ticker{
		interval: interval,
		jitter:   jitter,
		onTick:   onTick,
	}
}

// Run starts the ticker loop. It blocks until ctx is cancelled.
func (t *Ticker) Run(ctx context.Context) {
	for {
		wait := t.next()
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
			t.onTick(ctx)
		}
	}
}

// RunNow starts the ticker loop, firing the callback immediately before
// waiting for the first interval. It blocks until ctx is cancelled.
func (t *Ticker) RunNow(ctx context.Context) {
	t.onTick(ctx)
	t.Run(ctx)
}

// next returns the duration until the next tick, applying jitter if configured.
func (t *Ticker) next() time.Duration {
	if t.jitter == 0 {
		return t.interval
	}
	// Random offset in [-jitter, +jitter].
	offset := time.Duration(rand.Int63n(int64(t.jitter)*2+1)) - t.jitter
	return t.interval + offset
}
