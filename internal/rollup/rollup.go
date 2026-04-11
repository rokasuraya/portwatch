// Package rollup aggregates multiple port change events within a time
// window into a single summary, reducing alert noise during rapid churn.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Summary holds the net result of all events collected during a window.
type Summary struct {
	Opened  []snapshot.Entry
	Closed  []snapshot.Entry
	WindowEnd time.Time
}

// Rollup buffers port-change entries and flushes them after a quiet window.
type Rollup struct {
	mu      sync.Mutex
	opened  map[string]snapshot.Entry
	closed  map[string]snapshot.Entry
	timer   *time.Timer
	wait    time.Duration
	onFlush func(Summary)
}

// New creates a Rollup that calls onFlush after wait has elapsed with no
// new events. A wait of zero disables batching (flush is immediate).
func New(wait time.Duration, onFlush func(Summary)) *Rollup {
	return &Rollup{
		wait:    wait,
		opened:  make(map[string]snapshot.Entry),
		closed:  make(map[string]snapshot.Entry),
		onFlush: onFlush,
	}
}

// Add records opened and closed entries, resetting the flush timer.
func (r *Rollup) Add(opened, closed []snapshot.Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, e := range opened {
		r.opened[key(e)] = e
		delete(r.closed, key(e))
	}
	for _, e := range closed {
		r.closed[key(e)] = e
		delete(r.opened, key(e))
	}

	if r.wait == 0 {
		r.flush()
		return
	}

	if r.timer != nil {
		r.timer.Reset(r.wait)
	} else {
		r.timer = time.AfterFunc(r.wait, func() {
			r.mu.Lock()
			defer r.mu.Unlock()
			r.flush()
		})
	}
}

// Flush forces immediate emission of any buffered events.
func (r *Rollup) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.flush()
}

// flush emits the summary and resets internal state. Caller must hold mu.
func (r *Rollup) flush() {
	if len(r.opened) == 0 && len(r.closed) == 0 {
		return
	}
	s := Summary{WindowEnd: time.Now()}
	for _, e := range r.opened {
		s.Opened = append(s.Opened, e)
	}
	for _, e := range r.closed {
		s.Closed = append(s.Closed, e)
	}
	r.opened = make(map[string]snapshot.Entry)
	r.closed = make(map[string]snapshot.Entry)
	r.timer = nil
	r.onFlush(s)
}

func key(e snapshot.Entry) string { return e.Proto + ":" + e.String() }
