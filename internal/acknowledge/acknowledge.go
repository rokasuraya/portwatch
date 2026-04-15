// Package acknowledge provides a simple mechanism for tracking which
// port-change events have been acknowledged (silenced) by an operator.
// Acknowledged entries are persisted to disk so they survive restarts.
package acknowledge

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry records a single acknowledgement.
type Entry struct {
	Protocol  string    `json:"protocol"`
	Port      uint16    `json:"port"`
	AckedAt   time.Time `json:"acked_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// Acknowledger manages the set of acknowledged port events.
type Acknowledger struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
}

// New loads existing acknowledgements from path (if any) and returns an
// Acknowledger ready for use. A non-existent file is not an error.
func New(path string) (*Acknowledger, error) {
	a := &Acknowledger{
		entries: make(map[string]Entry),
		path:    path,
	}
	if err := a.load(); err != nil {
		return nil, err
	}
	return a, nil
}

// Ack acknowledges a port/protocol pair, optionally until a deadline.
// Passing a zero Time means the acknowledgement never expires.
func (a *Acknowledger) Ack(protocol string, port uint16, expires time.Time, note string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries[key(protocol, port)] = Entry{
		Protocol:  protocol,
		Port:      port,
		AckedAt:   time.Now(),
		ExpiresAt: expires,
		Note:      note,
	}
	return a.persist()
}

// IsAcked reports whether the port/protocol pair is currently acknowledged.
func (a *Acknowledger) IsAcked(protocol string, port uint16) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	e, ok := a.entries[key(protocol, port)]
	if !ok {
		return false
	}
	if !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt) {
		return false
	}
	return true
}

// Remove deletes an acknowledgement for the given port/protocol pair.
func (a *Acknowledger) Remove(protocol string, port uint16) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.entries, key(protocol, port))
	return a.persist()
}

// All returns a snapshot of all current (including expired) entries.
func (a *Acknowledger) All() []Entry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Entry, 0, len(a.entries))
	for _, e := range a.entries {
		out = append(out, e)
	}
	return out
}

func (a *Acknowledger) load() error {
	data, err := os.ReadFile(a.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("acknowledge: read %s: %w", a.path, err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("acknowledge: parse %s: %w", a.path, err)
	}
	for _, e := range entries {
		a.entries[key(e.Protocol, e.Port)] = e
	}
	return nil
}

func (a *Acknowledger) persist() error {
	list := make([]Entry, 0, len(a.entries))
	for _, e := range a.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("acknowledge: marshal: %w", err)
	}
	if err := os.WriteFile(a.path, data, 0o644); err != nil {
		return fmt.Errorf("acknowledge: write %s: %w", a.path, err)
	}
	return nil
}

func key(protocol string, port uint16) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
