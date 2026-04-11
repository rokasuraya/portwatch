package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Store persists the most recent Snapshot to disk and keeps the previous one
// in memory so callers can compute a Diff without extra bookkeeping.
type Store struct {
	mu       sync.RWMutex
	path     string
	current  Snapshot
	previous Snapshot
}

// NewStore creates a Store backed by the given file path.
// If the file already exists its contents are loaded as the current snapshot.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Set replaces the current snapshot, promoting the old one to previous,
// then persists the new snapshot to disk.
func (s *Store) Set(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.previous = s.current
	s.current = snap
	return s.persist()
}

// Current returns the most recently stored snapshot.
func (s *Store) Current() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Previous returns the snapshot that preceded the current one.
func (s *Store) Previous() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.previous
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.current)
}

func (s *Store) persist() error {
	data, err := json.Marshal(s.current)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
