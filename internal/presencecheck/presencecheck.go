// Package presencecheck verifies that expected ports are present in a snapshot
// and reports any that are absent, helping operators detect silent service failures.
package presencecheck

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Result holds the outcome of a single presence check.
type Result struct {
	Port     int
	Protocol string
	Present  bool
	CheckedAt time.Time
}

// Checker verifies that a set of required ports appear in a snapshot.
type Checker struct {
	required []snapshot.Entry
	out      io.Writer
}

// New returns a Checker that will verify the given required entries.
// If out is nil, os.Stdout is used.
func New(required []snapshot.Entry, out io.Writer) *Checker {
	if out == nil {
		out = os.Stdout
	}
	return &Checker{
		required: required,
		out: out,
	}
}

// Check compares required entries against the provided snapshot and returns
// a Result for each required entry.
func (c *Checker) Check(snap *snapshot.Snapshot) []Result {
	present := make(map[string]bool, len(snap.Entries))
	for _, e := range snap.Entries {
		present[key(e)] = true
	}

	now := time.Now()
	results := make([]Result, 0, len(c.required))
	for _, req := range c.required {
		results = append(results, Result{
			Port:      req.Port,
			Protocol:  req.Protocol,
			Present:   present[key(req)],
			CheckedAt: now,
		})
	}
	return results
}

// Report writes a human-readable summary of missing ports to the writer.
func (c *Checker) Report(results []Result) {
	for _, r := range results {
		if !r.Present {
			fmt.Fprintf(c.out, "[presencecheck] MISSING %s/%d at %s\n",
				r.Protocol, r.Port, r.CheckedAt.Format(time.RFC3339))
		}
	}
}

func key(e snapshot.Entry) string {
	return fmt.Sprintf("%s/%d", e.Protocol, e.Port)
}
