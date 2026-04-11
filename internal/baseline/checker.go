package baseline

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// CheckResult holds the outcome of a single baseline check.
type CheckResult struct {
	Timestamp  time.Time
	Unexpected []Entry
	Missing    []Entry
	Clean      bool
}

// Checker compares live snapshots against a Baseline and writes a human-
// readable summary to an io.Writer.
type Checker struct {
	baseline *Baseline
	out      io.Writer
}

// NewChecker creates a Checker that writes results to out.
// If out is nil, os.Stdout is used.
func NewChecker(b *Baseline, out io.Writer) *Checker {
	if out == nil {
		out = os.Stdout
	}
	return &Checker{baseline: b, out: out}
}

// Check compares snap against the baseline and returns a CheckResult.
// A summary is always written to the configured writer.
func (c *Checker) Check(snap *snapshot.Snapshot) CheckResult {
	unexpected, missing := c.baseline.Diff(snap)
	res := CheckResult{
		Timestamp:  time.Now().UTC(),
		Unexpected: unexpected,
		Missing:    missing,
		Clean:      len(unexpected) == 0 && len(missing) == 0,
	}
	c.write(res)
	return res
}

func (c *Checker) write(r CheckResult) {
	if r.Clean {
		fmt.Fprintf(c.out, "[%s] baseline check: OK\n", r.Timestamp.Format(time.RFC3339))
		return
	}
	fmt.Fprintf(c.out, "[%s] baseline check: DEVIATION\n", r.Timestamp.Format(time.RFC3339))
	for _, e := range r.Unexpected {
		fmt.Fprintf(c.out, "  UNEXPECTED  %s/%d\n", e.Protocol, e.Port)
	}
	for _, e := range r.Missing {
		fmt.Fprintf(c.out, "  MISSING     %s/%d\n", e.Protocol, e.Port)
	}
}
