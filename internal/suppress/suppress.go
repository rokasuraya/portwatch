// Package suppress provides a suppression list that prevents alerts
// for ports that have been explicitly silenced by the operator.
package suppress

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Entry describes a single supp.
type Entry struct {
	Port     int
	Protocol string
	Reason   string
	ExpiresAt time.Time // zero means never expires
}

// Suppressor holds the active suppression list.
type Suppressor struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Suppressor.
func New() *Suppressor {
	return &Suppressor{
		entries: make(map[string]Entry),
	}
}

// Add registers a suppression rule. If expiresAt is zero the rule never expires.
func (s *Suppressor) Add(port int, protocol, reason string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key(port, protocol)] = Entry{
		Port:      port,
		Protocol:  protocol,
		Reason:    reason,
		ExpiresAt: expiresAt,
	}
}

// Remove deletes a suppression rule.
func (s *Suppressor) Remove(port int, protocol string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key(port, protocol))
}

// IsSuppressed reports whether the given port/protocol combination is
// currently suppressed. Expired rules are treated as absent and are
// lazily removed.
func (s *Suppressor) IsSuppressed(port int, protocol string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key(port, protocol)]
	if !ok {
		return false
	}
	if !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt) {
		delete(s.entries, key(port, protocol))
		return false
	}
	return true
}

// List returns a snapshot of all active (non-expired) suppression entries.
func (s *Suppressor) List() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	out := make([]Entry, 0, len(s.entries))
	for k, e := range s.entries {
		if !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt) {
			delete(s.entries, k)
			continue
		}
		out = append(out, e)
	}
	return out
}

func key(port int, protocol string) string {
	return fmt.Sprintf("%d/%s", port, strings.ToLower(protocol))
}
