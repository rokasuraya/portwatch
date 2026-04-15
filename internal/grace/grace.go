// Package grace provides a graceful shutdown coordinator that waits for
// in-flight work to complete before allowing the process to exit.
package grace

import (
	"context"
	"sync"
	"time"
)

// Coordinator tracks active workers and blocks shutdown until all finish
// or the deadline is exceeded.
type Coordinator struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	timeout time.Duration
	done    chan struct{}
}

// New returns a Coordinator with the given shutdown timeout.
func New(timeout time.Duration) *Coordinator {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Coordinator{
		timeout: timeout,
		done:    make(chan struct{}),
	}
}

// Acquire registers one unit of in-flight work. It returns false if the
// coordinator has already begun shutting down.
func (c *Coordinator) Acquire() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.done:
		return false
	default:
	}
	c.wg.Add(1)
	return true
}

// Release marks one unit of in-flight work as complete.
func (c *Coordinator) Release() {
	c.wg.Done()
}

// Shutdown signals that no new work should be accepted and waits for all
// in-flight workers to finish. It returns ctx.Err() if the parent context
// is cancelled before workers drain, or context.DeadlineExceeded if the
// internal timeout fires first.
func (c *Coordinator) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	select {
	case <-c.done:
		c.mu.Unlock()
		return nil
	default:
		close(c.done)
	}
	c.mu.Unlock()

	finished := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(finished)
	}()

	deadline := time.NewTimer(c.timeout)
	defer deadline.Stop()

	select {
	case <-finished:
		return nil
	case <-deadline.C:
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}
