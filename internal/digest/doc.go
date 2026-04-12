// Package digest provides snapshot fingerprinting for portwatch.
//
// A Digest is a stable SHA-256 hash of the sorted set of open (proto, port)
// pairs observed during a scan. Digests allow other components to skip
// expensive diffing and alerting work when nothing has changed between two
// consecutive scans.
//
// Typical usage:
//
//	tracker := digest.New()
//	// ... inside scan loop:
//	d, changed := tracker.Changed(snap.Entries)
//	if !changed {
//	    return // nothing to do
//	}
//	log.Printf("state changed, new digest: %s", d)
package digest
