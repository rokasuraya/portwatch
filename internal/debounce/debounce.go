// Package debounce provides a simple debouncer that suppresses rapid
// repeated triggers and only fires after a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays execution of a function until after a specified wait
// duration has passed since the last call to Trigger.
type Debouncer struct {
	wait  time.Duration
	mu    sync.Mutex
	timer *time.Timer
}

// New creates a new Debouncer with the given wait duration.
func New(wait time.Duration) *Debouncer {
	return &Debouncer{wait: wait}
}

// Trigger schedules fn to be called after the wait duration. If Trigger is
// called again before the timer fires, the timer resets and fn is delayed
// further. Only the most recent fn is executed.
func (d *Debouncer) Trigger(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.wait, fn)
}

// Flush cancels any pending timer and executes fn immediately if one was
// scheduled. Returns true if a pending call was flushed.
func (d *Debouncer) Flush() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer == nil {
		return false
	}
	stopped := d.timer.Stop()
	d.timer = nil
	return stopped
}

// Reset cancels any pending timer without executing fn.
func (d *Debouncer) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
