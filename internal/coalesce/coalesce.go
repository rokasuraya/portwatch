// Package coalesce batches port-change diff entries that arrive within a
// short window and emits them as a single combined slice.
package coalesce

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Coalescer accumulates diff entries and flushes them after a quiet window.
type Coalescer struct {
	mu      sync.Mutex
	wait    time.Duration
	onFlush func(opened, closed []snapshot.Entry)
	opened  []snapshot.Entry
	closed  []snapshot.Entry
	timer   *time.Timer
}

// New returns a Coalescer that waits for wait after the last call to Add
// before invoking onFlush with the accumulated entries.
func New(wait time.Duration, onFlush func(opened, closed []snapshot.Entry)) *Coalescer {
	if wait <= 0 {
		wait = 200 * time.Millisecond
	}
	return &Coalescer{wait: wait, onFlush: onFlush}
}

// Add appends entries to the pending batch and resets the flush timer.
func (c *Coalescer) Add(opened, closed []snapshot.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.opened = append(c.opened, opened...)
	c.closed = append(c.closed, closed...)

	if c.timer != nil {
		c.timer.Stop()
	}
	c.timer = time.AfterFunc(c.wait, c.flush)
}

// Flush cancels any pending timer and immediately emits accumulated entries.
// Returns true if there were pending entries to flush.
func (c *Coalescer) Flush() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	return c.emit()
}

// flush is called by the internal timer (no lock held by caller).
func (c *Coalescer) flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.emit()
}

// emit fires the callback and resets buffers. Caller must hold mu.
func (c *Coalescer) emit() bool {
	if len(c.opened) == 0 && len(c.closed) == 0 {
		return false
	}
	opened := c.opened
	closed := c.closed
	c.opened = nil
	c.closed = nil
	c.timer = nil
	c.onFlush(opened, closed)
	return true
}
