// Package dedup provides a scan-result deduplicator that suppresses
// identical consecutive snapshots so downstream consumers only receive
// meaningful changes.
package dedup

import (
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Deduplicator tracks the fingerprint of the last accepted snapshot and
// drops any snapshot whose fingerprint matches the previous one.
type Deduplicator struct {
	mu   sync.Mutex
	last string
}

// New returns a ready-to-use Deduplicator.
func New() *Deduplicator {
	return &Deduplicator{}
}

// Accept returns true if snap is different from the previously accepted
// snapshot, and records it as the new baseline. It returns false when the
// snapshot is identical to the last one seen.
func (d *Deduplicator) Accept(snap *snapshot.Snapshot) bool {
	if snap == nil {
		return false
	}

	fp := fingerprint(snap)

	d.mu.Lock()
	defer d.mu.Unlock()

	if fp == d.last {
		return false
	}

	d.last = fp
	return true
}

// Reset clears the stored fingerprint so the next call to Accept always
// returns true regardless of content.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = ""
}

// fingerprint builds a lightweight string key from the snapshot entries.
// Order-independent so port set {80,443} == {443,80}.
func fingerprint(snap *snapshot.Snapshot) string {
	entries := snap.Entries()
	if len(entries) == 0 {
		return "empty"
	}

	// Use a simple sorted concatenation for the key.
	seen := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		seen[e.String()] = struct{}{}
	}

	var buf []byte
	for k := range seen {
		buf = append(buf, k...)
		buf = append(buf, '|')
	}
	return string(buf)
}
