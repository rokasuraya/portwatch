package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int       `json:"port"`
	Protocol string    `json:"protocol"`
	Open     bool      `json:"open"`
	SeenAt   time.Time `json:"seen_at"`
}

// Snapshot holds a collection of port states at a point in time.
type Snapshot struct {
	Timestamp time.Time             `json:"timestamp"`
	Ports     map[int]*PortState    `json:"ports"`
}

// Store manages port state persistence and change detection.
type Store struct {
	mu       sync.RWMutex
	current  *Snapshot
	filePath string
}

// New creates a new Store, loading existing state from filePath if present.
func New(filePath string) (*Store, error) {
	s := &Store{
		filePath: filePath,
		current: &Snapshot{
			Ports: make(map[int]*PortState),
		},
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Diff represents a change between two snapshots.
type Diff struct {
	Opened []PortState
	Closed []PortState
}

// Update compares newPorts against the current snapshot, persists the new
// state, and returns a Diff describing what changed.
func (s *Store) Update(newPorts []PortState) (Diff, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newMap := make(map[int]*PortState, len(newPorts))
	for i := range newPorts {
		p := newPorts[i]
		newMap[p.Port] = &p
	}

	var diff Diff
	for port, ps := range newMap {
		if _, existed := s.current.Ports[port]; !existed {
			diff.Opened = append(diff.Opened, *ps)
		}
	}
	for port, ps := range s.current.Ports {
		if _, stillOpen := newMap[port]; !stillOpen {
			diff.Closed = append(diff.Closed, *ps)
		}
	}

	s.current = &Snapshot{
		Timestamp: time.Now(),
		Ports:     newMap,
	}
	return diff, s.save()
}

// Current returns a copy of the current snapshot.
func (s *Store) Current() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := *s.current
	return copy
}

func (s *Store) load() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(s.current)
}

func (s *Store) save() error {
	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s.current)
}
