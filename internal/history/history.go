// Package history maintains a rolling log of port change events
// so the CLI can display recent activity without re-scanning.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Event represents a single port-change occurrence.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      string    `json:"kind"` // "opened" | "closed"
	Proto     string    `json:"proto"`
	Port      int       `json:"port"`
}

// History holds a capped, ordered list of events.
type History struct {
	mu     sync.Mutex
	events []Event
	cap    int
	path   string
}

// New creates a History that keeps at most maxEvents entries and
// persists them to path. Existing entries are loaded automatically.
func New(path string, maxEvents int) (*History, error) {
	h := &History{path: path, cap: maxEvents}
	if err := h.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return h, nil
}

// Record appends a new event, evicting the oldest when over capacity.
func (h *History) Record(kind, proto string, port int) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Proto:     proto,
		Port:      port,
	})
	if len(h.events) > h.cap {
		h.events = h.events[len(h.events)-h.cap:]
	}
	return h.persist()
}

// Events returns a copy of all stored events, oldest first.
func (h *History) Events() []Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Event, len(h.events))
	copy(out, h.events)
	return out
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.events)
}

func (h *History) persist() error {
	data, err := json.MarshalIndent(h.events, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o600)
}
