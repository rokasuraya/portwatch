// Package stale detects ports that have remained open beyond a configured
// duration and emits warnings so operators can review long-lived listeners.
package stale

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry records when a port was first seen open.
type Entry struct {
	Proto   string
	Port    int
	FirstAt time.Time
}

// Detector tracks open ports and warns when they exceed the max age.
type Detector struct {
	mu     sync.Mutex
	seen   map[string]Entry
	maxAge time.Duration
	out    io.Writer
}

// New returns a Detector that warns about ports open longer than maxAge.
// Warnings are written to out; if out is nil, os.Stderr is used.
func New(maxAge time.Duration, out io.Writer) *Detector {
	if out == nil {
		out = os.Stderr
	}
	return &Detector{
		seen:   make(map[string]Entry),
		maxAge: maxAge,
		out:    out,
	}
}

func portKey(proto string, port int) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Observe updates internal state from the current snapshot.
// New ports are recorded; ports no longer present are removed.
func (d *Detector) Observe(snap *snapshot.Snapshot) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	active := make(map[string]struct{}, len(snap.Entries))

	for _, e := range snap.Entries {
		k := portKey(e.Proto, e.Port)
		active[k] = struct{}{}
		if _, ok := d.seen[k]; !ok {
			d.seen[k] = Entry{Proto: e.Proto, Port: e.Port, FirstAt: now}
		}
	}

	for k := range d.seen {
		if _, ok := active[k]; !ok {
			delete(d.seen, k)
		}
	}
}

// Check emits a warning for every tracked port whose age exceeds maxAge.
// It returns the number of stale ports found.
func (d *Detector) Check() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	count := 0
	for _, e := range d.seen {
		age := now.Sub(e.FirstAt)
		if age >= d.maxAge {
			fmt.Fprintf(d.out, "[stale] %s port %d open for %s (exceeds %s)\n",
				e.Proto, e.Port, age.Round(time.Second), d.maxAge)
			count++
		}
	}
	return count
}
