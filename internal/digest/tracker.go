package digest

import (
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Tracker keeps the last-seen digest and reports whether a new snapshot
// differs from the previous one.
type Tracker struct {
	mu   sync.Mutex
	last Digest
}

// New returns a new Tracker with no prior digest.
func New() *Tracker {
	return &Tracker{}
}

// Changed computes the digest of snap and returns (digest, true) when it
// differs from the previously recorded digest, or (digest, false) when the
// snapshot is unchanged. The internal state is always updated to the new
// digest.
func (t *Tracker) Changed(entries []snapshot.Entry) (Digest, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	d := Compute(entries)
	changed := d != t.last
	t.last = d
	return d, changed
}

// Last returns the most recently recorded digest.
func (t *Tracker) Last() Digest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.last
}

// Reset clears the stored digest so the next call to Changed always
// reports a change.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = ""
}
