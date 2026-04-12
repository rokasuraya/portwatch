// Package digest computes and compares fingerprints of port snapshots.
// A digest is a stable hash of the sorted set of open ports, allowing
// callers to cheaply detect whether a scan result has changed.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/snapshot"
)

// Digest is a hex-encoded SHA-256 fingerprint of a snapshot's entries.
type Digest string

// Compute returns a Digest for the given snapshot entries.
// The entries are sorted by protocol+port before hashing so the result
// is independent of scan order.
func Compute(entries []snapshot.Entry) Digest {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = fmt.Sprintf("%s:%d", e.Proto, e.Port)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintln(h, k)
	}
	return Digest(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool {
	return a == b
}

// String returns the hex string representation of the digest.
func (d Digest) String() string {
	return string(d)
}
