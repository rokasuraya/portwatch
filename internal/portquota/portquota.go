// Package portquota enforces a maximum number of concurrently open ports
// and emits a warning when the threshold is exceeded.
package portquota

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Quota holds the configuration and state for port quota enforcement.
type Quota struct {
	mu        sync.Mutex
	max       int
	out       io.Writer
	exceeded  bool
}

// New returns a Quota that warns when open port count exceeds max.
// If out is nil it defaults to os.Stderr.
func New(max int, out io.Writer) *Quota {
	if out == nil {
		out = os.Stderr
	}
	return &Quota{max: max, out: out}
}

// Check evaluates the snapshot against the quota limit.
// It writes a warning to the configured writer when the limit is exceeded
// and returns true if the quota was breached.
func (q *Quota) Check(snap *snapshot.Snapshot) bool {
	if snap == nil {
		return false
	}
	count := len(snap.Entries)
	q.mu.Lock()
	defer q.mu.Unlock()
	if count > q.max {
		if !q.exceeded {
			fmt.Fprintf(q.out, "portquota: open port count %d exceeds limit %d\n", count, q.max)
		}
		q.exceeded = true
		return true
	}
	q.exceeded = false
	return false
}

// Max returns the configured quota limit.
func (q *Quota) Max() int {
	return q.max
}

// SetMax updates the quota limit.
func (q *Quota) SetMax(max int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.max = max
}
