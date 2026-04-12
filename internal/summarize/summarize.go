// Package summarize provides a periodic summary reporter that aggregates
// scan metrics and emits a human-readable digest at a configurable interval.
package summarize

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of scan activity.
type Snapshot struct {
	At           time.Time
	TotalScans   int
	TotalOpened  int
	TotalClosed  int
	TotalAlerts  int
}

// Summarizer collects incremental scan data and writes a digest
// to an io.Writer at a fixed interval.
type Summarizer struct {
	mu       sync.Mutex
	current  Snapshot
	interval time.Duration
	out      io.Writer
}

// New returns a Summarizer that flushes to out every interval.
// If out is nil, os.Stdout is used.
func New(interval time.Duration, out io.Writer) *Summarizer {
	if out == nil {
		out = os.Stdout
	}
	return &Summarizer{
		interval: interval,
		out:      out,
	}
}

// Record accumulates counts from a single scan tick.
func (s *Summarizer) Record(opened, closed, alerts int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current.TotalScans++
	s.current.TotalOpened += opened
	s.current.TotalClosed += closed
	s.current.TotalAlerts += alerts
}

// flush writes the current snapshot to the writer and resets counters.
func (s *Summarizer) flush() {
	s.mu.Lock()
	snap := s.current
	snap.At = time.Now()
	s.current = Snapshot{}
	s.mu.Unlock()

	fmt.Fprintf(s.out,
		"[summary] %s | scans=%d opened=%d closed=%d alerts=%d\n",
		snap.At.Format(time.RFC3339),
		snap.TotalScans,
		snap.TotalOpened,
		snap.TotalClosed,
		snap.TotalAlerts,
	)
}

// Run starts the summary flush loop. It blocks until ctx is cancelled.
func (s *Summarizer) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.flush()
			return
		case <-ticker.C:
			s.flush()
		}
	}
}
