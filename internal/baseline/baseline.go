// Package baseline manages the expected (trusted) set of open ports.
// A baseline represents the operator-approved port state; deviations from
// it are surfaced as alerts.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry is a single approved port/protocol pair.
type Entry struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// Baseline holds the approved set of entries and the time it was last updated.
type Baseline struct {
	mu      sync.RWMutex
	path    string
	Entries []Entry    `json:"entries"`
	Updated time.Time  `json:"updated"`
}

// New loads a baseline from path, or returns an empty baseline if the file
// does not exist yet.
func New(path string) (*Baseline, error) {
	b := &Baseline{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, b); err != nil {
		return nil, err
	}
	return b, nil
}

// Approve replaces the current baseline with the entries from snap and
// persists the result to disk.
func (b *Baseline) Approve(snap *snapshot.Snapshot) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	entries := make([]Entry, 0, len(snap.Entries))
	for _, e := range snap.Entries {
		entries = append(entries, Entry{Port: e.Port, Protocol: e.Protocol})
	}
	b.Entries = entries
	b.Updated = time.Now().UTC()
	return b.save()
}

// Diff returns the entries in snap that are not present in the baseline
// (unexpected ports) and the entries in the baseline that are absent from
// snap (missing ports).
func (b *Baseline) Diff(snap *snapshot.Snapshot) (unexpected []Entry, missing []Entry) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	approved := make(map[string]struct{}, len(b.Entries))
	for _, e := range b.Entries {
		approved[key(e.Port, e.Protocol)] = struct{}{}
	}

	seen := make(map[string]struct{}, len(snap.Entries))
	for _, e := range snap.Entries {
		k := key(e.Port, e.Protocol)
		seen[k] = struct{}{}
		if _, ok := approved[k]; !ok {
			unexpected = append(unexpected, Entry{Port: e.Port, Protocol: e.Protocol})
		}
	}

	for _, e := range b.Entries {
		if _, ok := seen[key(e.Port, e.Protocol)]; !ok {
			missing = append(missing, e)
		}
	}
	return
}

func (b *Baseline) save() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

func key(port int, protocol string) string {
	return protocol + ":" + string(rune(port))
}
