// Package portlock provides a mechanism to lock (pin) specific ports so that
// changes to them are silently ignored by the alerting pipeline.
package portlock

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Entry identifies a locked port.
type Entry struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Reason   string `json:"reason,omitempty"`
}

// Locker holds the set of locked ports.
type Locker struct {
	mu      sync.RWMutex
	locked  map[string]Entry
	output  io.Writer
}

func key(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// New returns a new Locker. If w is nil, os.Stderr is used.
func New(w io.Writer) *Locker {
	if w == nil {
		w = os.Stderr
	}
	return &Locker{
		locked: make(map[string]Entry),
		output: w,
	}
}

// Lock adds a port to the locked set.
func (l *Locker) Lock(port int, protocol, reason string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locked[key(port, protocol)] = Entry{Port: port, Protocol: protocol, Reason: reason}
}

// Unlock removes a port from the locked set.
func (l *Locker) Unlock(port int, protocol string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.locked, key(port, protocol))
}

// IsLocked reports whether the given port/protocol pair is locked.
func (l *Locker) IsLocked(port int, protocol string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.locked[key(port, protocol)]
	return ok
}

// Len returns the number of currently locked ports.
func (l *Locker) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.locked)
}

// Entries returns a snapshot of all locked entries.
func (l *Locker) Entries() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0, len(l.locked))
	for _, e := range l.locked {
		out = append(out, e)
	}
	return out
}
