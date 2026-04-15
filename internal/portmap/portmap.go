// Package portmap maintains a live mapping of port numbers to their
// current open/closed status, enriched with a human-readable label.
package portmap

import (
	"fmt"
	"sync"
)

// Entry holds the current state of a single port.
type Entry struct {
	Port     int
	Protocol string
	Label    string
	Open     bool
}

// PortMap is a thread-safe registry of port entries.
type PortMap struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised, empty PortMap.
func New() *PortMap {
	return &PortMap{
		entries: make(map[string]*Entry),
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Set inserts or updates the entry for the given port/protocol pair.
func (pm *PortMap) Set(port int, proto, label string, open bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.entries[key(port, proto)] = &Entry{
		Port:     port,
		Protocol: proto,
		Label:    label,
		Open:     open,
	}
}

// Get returns the entry for the given port/protocol pair and whether it exists.
func (pm *PortMap) Get(port int, proto string) (*Entry, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	e, ok := pm.entries[key(port, proto)]
	if !ok {
		return nil, false
	}
	copy := *e
	return &copy, true
}

// Delete removes the entry for the given port/protocol pair.
func (pm *PortMap) Delete(port int, proto string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.entries, key(port, proto))
}

// Len returns the number of tracked entries.
func (pm *PortMap) Len() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.entries)
}

// All returns a snapshot of all current entries.
func (pm *PortMap) All() []Entry {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	out := make([]Entry, 0, len(pm.entries))
	for _, e := range pm.entries {
		out = append(out, *e)
	}
	return out
}
