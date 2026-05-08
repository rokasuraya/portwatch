// Package portschedule provides a scan schedule registry that maps port ranges
// to custom scan intervals, allowing high-risk ports to be scanned more
// frequently than low-priority ranges.
package portschedule

import (
	"fmt"
	"sync"
	"time"
)

// Entry associates a port range with a desired scan interval.
type Entry struct {
	MinPort  int
	MaxPort  int
	Interval time.Duration
	Label    string
}

// Schedule holds a set of port range entries and resolves the shortest
// matching interval for a given port number.
type Schedule struct {
	mu      sync.RWMutex
	entries []Entry
	default_ time.Duration
}

// New returns a Schedule with the given default interval used when no entry
// matches a queried port.
func New(defaultInterval time.Duration) *Schedule {
	if defaultInterval <= 0 {
		defaultInterval = 60 * time.Second
	}
	return &Schedule{default_: defaultInterval}
}

// Add registers a port range entry. Returns an error if the range is invalid.
func (s *Schedule) Add(e Entry) error {
	if e.MinPort < 0 || e.MaxPort > 65535 || e.MinPort > e.MaxPort {
		return fmt.Errorf("portschedule: invalid range %d-%d", e.MinPort, e.MaxPort)
	}
	if e.Interval <= 0 {
		return fmt.Errorf("portschedule: interval must be positive")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	return nil
}

// Interval returns the shortest matching interval for port, or the default
// interval if no entry covers it.
func (s *Schedule) Interval(port int) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	best := s.default_
	for _, e := range s.entries {
		if port >= e.MinPort && port <= e.MaxPort {
			if e.Interval < best {
				best = e.Interval
			}
		}
	}
	return best
}

// Len returns the number of registered entries.
func (s *Schedule) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// Reset removes all entries, leaving the default interval intact.
func (s *Schedule) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = s.entries[:0]
}
