// Package snapshot provides point-in-time captures of open port state
// along with utilities to diff two snapshots.
package snapshot

import (
	"fmt"
	"time"
)

// Entry represents a single open port at a point in time.
type Entry struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

// String returns a human-readable representation of the entry.
func (e Entry) String() string {
	return fmt.Sprintf("%s:%d", e.Protocol, e.Port)
}

// Snapshot is an immutable capture of open ports at a given time.
type Snapshot struct {
	CapturedAt time.Time
	Entries    []Entry
}

// New creates a new Snapshot from the provided entries, stamped with now.
func New(entries []Entry) Snapshot {
	copied := make([]Entry, len(entries))
	copy(copied, entries)
	return Snapshot{
		CapturedAt: time.Now(),
		Entries:    copied,
	}
}

// Diff holds the ports that appeared or disappeared between two snapshots.
type Diff struct {
	Opened []Entry
	Closed []Entry
}

// IsEmpty returns true when there are no changes between the two snapshots.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Compare returns a Diff between a previous and current Snapshot.
// Entries present in current but not previous are Opened;
// entries present in previous but not current are Closed.
func Compare(previous, current Snapshot) Diff {
	prev := index(previous.Entries)
	curr := index(current.Entries)

	var opened, closed []Entry

	for key, e := range curr {
		if _, ok := prev[key]; !ok {
			opened = append(opened, e)
		}
	}
	for key, e := range prev {
		if _, ok := curr[key]; !ok {
			closed = append(closed, e)
		}
	}

	return Diff{Opened: opened, Closed: closed}
}

func index(entries []Entry) map[string]Entry {
	m := make(map[string]Entry, len(entries))
	for _, e := range entries {
		m[fmt.Sprintf("%s:%d", e.Protocol, e.Port)] = e
	}
	return m
}
