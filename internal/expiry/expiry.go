// Package expiry tracks port entries and flags those that have been
// continuously open beyond a configurable maximum age.
package expiry

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
	Port     int
	Protocol string
	FirstSeen time.Time
}

// Checker tracks first-seen timestamps for open ports and writes warnings
// for entries that exceed MaxAge.
type Checker struct {
	mu      sync.Mutex
	records map[string]Entry
	MaxAge  time.Duration
	out     io.Writer
}

// New returns a Checker that warns when a port has been open longer than maxAge.
// Output is written to w; pass nil to default to os.Stdout.
func New(maxAge time.Duration, w io.Writer) *Checker {
	if w == nil {
		w = os.Stdout
	}
	return &Checker{
		records: make(map[string]Entry),
		MaxAge:  maxAge,
		out:     w,
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Observe updates internal state from the current snapshot, recording
// first-seen times for new ports and removing entries that are no longer open.
func (c *Checker) Observe(snap *snapshot.Snapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()

	seen := make(map[string]struct{}, len(snap.Entries))
	for _, e := range snap.Entries {
		k := portKey(e.Port, e.Protocol)
		seen[k] = struct{}{}
		if _, exists := c.records[k]; !exists {
			c.records[k] = Entry{
				Port:      e.Port,
				Protocol:  e.Protocol,
				FirstSeen: time.Now(),
			}
		}
	}

	for k := range c.records {
		if _, ok := seen[k]; !ok {
			delete(c.records, k)
		}
	}
}

// Check writes a warning line to the configured writer for every tracked port
// whose age exceeds MaxAge. It returns the number of expired entries found.
func (c *Checker) Check() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	count := 0
	for _, rec := range c.records {
		age := now.Sub(rec.FirstSeen)
		if age > c.MaxAge {
			fmt.Fprintf(c.out, "[expiry] %s port %d open for %s (exceeds %s)\n",
				rec.Protocol, rec.Port, age.Round(time.Second), c.MaxAge)
			count++
		}
	}
	return count
}
