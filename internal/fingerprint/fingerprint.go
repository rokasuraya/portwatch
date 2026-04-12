// Package fingerprint provides port-state fingerprinting to detect
// meaningful changes between successive scans by hashing the set of
// open entries and comparing against a stored reference.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a hex-encoded SHA-256 hash of a sorted set of port entries.
type Fingerprint string

// Tracker stores the last observed fingerprint and reports whether a new
// snapshot represents a change.
type Tracker struct {
	mu   sync.Mutex
	last Fingerprint
}

// New returns a zero-value Tracker with no prior fingerprint.
func New() *Tracker {
	return &Tracker{}
}

// Compute derives a deterministic Fingerprint from a slice of scanner entries.
// Ordering of the input slice does not affect the result.
func Compute(entries []scanner.Entry) Fingerprint {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = fmt.Sprintf("%s:%d", e.Protocol, e.Port)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintln(h, k)
	}
	return Fingerprint(hex.EncodeToString(h.Sum(nil)))
}

// Changed returns true when fp differs from the last stored fingerprint and
// atomically updates the stored value to fp.
func (t *Tracker) Changed(fp Fingerprint) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if fp == t.last {
		return false
	}
	t.last = fp
	return true
}

// Last returns the most recently stored fingerprint.
func (t *Tracker) Last() Fingerprint {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.last
}

// Reset clears the stored fingerprint so the next call to Changed always
// returns true.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = ""
}
