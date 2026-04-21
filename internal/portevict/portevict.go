// Package portevict tracks ports that have been evicted (forcibly removed
// from the monitored set) and prevents them from re-triggering alerts
// until a configurable quiet period has elapsed.
package portevict

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Evictor records eviction timestamps and answers whether a port is
// still within its quiet period.
type Evictor struct {
	mu      sync.Mutex
	evicted map[string]time.Time
	quiet   time.Duration
	now     func() time.Time
}

// New returns an Evictor with the given quiet period.
func New(quiet time.Duration) *Evictor {
	return &Evictor{
		evicted: make(map[string]time.Time),
		quiet:   quiet,
		now:     time.Now,
	}
}

// Evict marks the entry as evicted at the current time.
func (e *Evictor) Evict(entry snapshot.Entry) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.evicted[portKey(entry)] = e.now()
}

// IsEvicted reports whether the entry is still within its quiet period.
func (e *Evictor) IsEvicted(entry snapshot.Entry) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	at, ok := e.evicted[portKey(entry)]
	if !ok {
		return false
	}
	if e.now().Sub(at) >= e.quiet {
		delete(e.evicted, portKey(entry))
		return false
	}
	return true
}

// Clear removes the eviction record for the given entry immediately.
func (e *Evictor) Clear(entry snapshot.Entry) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.evicted, portKey(entry))
}

// Len returns the number of currently evicted ports.
func (e *Evictor) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.evicted)
}

func portKey(entry snapshot.Entry) string {
	return fmt.Sprintf("%d/%s", entry.Port, entry.Protocol)
}
