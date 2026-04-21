// Package portdiff computes a human-readable diff summary between two
// snapshots, annotating each entry with its label and classification.
package portdiff

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry represents a single line in the diff output.
type Entry struct {
	Op       string // "opened" or "closed"
	Port     int
	Protocol string
	Label    string
}

// String returns a short human-readable representation.
func (e Entry) String() string {
	if e.Label != "" {
		return fmt.Sprintf("%s %s/%d (%s)", e.Op, e.Protocol, e.Port, e.Label)
	}
	return fmt.Sprintf("%s %s/%d", e.Op, e.Protocol, e.Port)
}

// Labeler maps a port+protocol pair to a descriptive name.
type Labeler interface {
	Label(port int, proto string) string
}

// Diff holds the result of comparing two snapshots.
type Diff struct {
	Opened []Entry
	Closed []Entry
}

// IsEmpty reports whether there are no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Summary returns a one-line description of the diff.
func (d Diff) Summary() string {
	parts := make([]string, 0, 2)
	if n := len(d.Opened); n > 0 {
		parts = append(parts, fmt.Sprintf("%d opened", n))
	}
	if n := len(d.Closed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d closed", n))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

// Compute derives a Diff between prev and next, annotating entries via l.
// Either snapshot may be nil; a nil prev treats all next entries as opened.
func Compute(prev, next *snapshot.Snapshot, l Labeler) Diff {
	opened, closed := snapshot.Compare(prev, next)

	toEntries := func(raw []snapshot.Entry, op string) []Entry {
		out := make([]Entry, 0, len(raw))
		for _, e := range raw {
			label := ""
			if l != nil {
				label = l.Label(e.Port, e.Protocol)
			}
			out = append(out, Entry{
				Op:       op,
				Port:     e.Port,
				Protocol: e.Protocol,
				Label:    label,
			})
		}
		return out
	}

	return Diff{
		Opened: toEntries(opened, "opened"),
		Closed: toEntries(closed, "closed"),
	}
}
