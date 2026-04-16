// Package gwatcher watches for port group membership changes and emits
// notifications when a port transitions into or out of a named group.
package gwatcher

import (
	"fmt"
	"io"
	"os"
	"sync"

	"portwatch/internal/portgroup"
	"portwatch/internal/snapshot"
)

// Event describes a group membership change for a single port.
type Event struct {
	Port     int
	Proto    string
	Group    string
	Joined   bool // true = joined, false = left
}

// Watcher emits Events when ports join or leave groups.
type Watcher struct {
	mu      sync.Mutex
	matcher *portgroup.Matcher
	prev    map[string][]string // key(port,proto) -> groups
	out     io.Writer
}

// New returns a Watcher backed by the given Matcher.
// Notifications are written to w; if nil, os.Stdout is used.
func New(m *portgroup.Matcher, w io.Writer) *Watcher {
	if w == nil {
		w = os.Stdout
	}
	return &Watcher{
		matcher: m,
		prev:    make(map[string][]string),
		out:     w,
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Observe compares the current snapshot against the previous one and emits
// group-change events for any port whose group membership has changed.
func (w *Watcher) Observe(snap *snapshot.Snapshot) []Event {
	if snap == nil {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	var events []Event
	current := make(map[string][]string)

	for _, e := range snap.Entries {
		k := portKey(e.Port, e.Proto)
		groups := w.matcher.Match(e.Port, e.Proto)
		current[k] = groups

		prevGroups := toSet(w.prev[k])
		for _, g := range groups {
			if !prevGroups[g] {
				events = append(events, Event{Port: e.Port, Proto: e.Proto, Group: g, Joined: true})
				fmt.Fprintf(w.out, "port %d/%s joined group %q\n", e.Port, e.Proto, g)
			}
		}
		currSet := toSet(groups)
		for _, g := range w.prev[k] {
			if !currSet[g] {
				events = append(events, Event{Port: e.Port, Proto: e.Proto, Group: g, Joined: false})
				fmt.Fprintf(w.out, "port %d/%s left group %q\n", e.Port, e.Proto, g)
			}
		}
	}
	w.prev = current
	return events
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
