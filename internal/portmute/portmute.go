// Package portmute provides a temporary muting mechanism for specific ports,
// suppressing alerts for a configurable duration without permanently ignoring them.
package portmute

import (
	"fmt"
	"sync"
	"time"
)

// Mute represents an active mute rule for a port/protocol pair.
type Mute struct {
	Port     int
	Protocol string
	Until    time.Time
	Reason   string
}

// Muter holds active mute rules and suppresses alerts for muted ports.
type Muter struct {
	mu    sync.Mutex
	rules map[string]Mute
	now   func() time.Time
}

// New returns a new Muter with no active rules.
func New() *Muter {
	return &Muter{
		rules: make(map[string]Mute),
		now:   time.Now,
	}
}

func portKey(port int, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}

// Mute suppresses alerts for the given port/protocol for the specified duration.
func (m *Muter) Mute(port int, protocol string, duration time.Duration, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules[portKey(port, protocol)] = Mute{
		Port:     port,
		Protocol: protocol,
		Until:    m.now().Add(duration),
		Reason:   reason,
	}
}

// Unmute removes a mute rule for the given port/protocol pair.
func (m *Muter) Unmute(port int, protocol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rules, portKey(port, protocol))
}

// IsMuted returns true if the port/protocol pair is currently muted.
func (m *Muter) IsMuted(port int, protocol string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	rule, ok := m.rules[portKey(port, protocol)]
	if !ok {
		return false
	}
	if m.now().After(rule.Until) {
		delete(m.rules, portKey(port, protocol))
		return false
	}
	return true
}

// Active returns a snapshot of all currently active mute rules.
func (m *Muter) Active() []Mute {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	out := make([]Mute, 0, len(m.rules))
	for k, rule := range m.rules {
		if now.After(rule.Until) {
			delete(m.rules, k)
			continue
		}
		out = append(out, rule)
	}
	return out
}
