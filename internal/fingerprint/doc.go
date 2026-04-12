// Package fingerprint provides lightweight port-state fingerprinting for
// portwatch.
//
// A Fingerprint is a deterministic SHA-256 hash derived from the set of open
// port entries observed during a scan. Because the hash is order-independent,
// two scans that find the same ports in a different order produce identical
// fingerprints.
//
// The Tracker type wraps a single stored fingerprint and exposes Changed, which
// atomically compares a new fingerprint against the stored one and updates it
// only when a difference is detected. This makes it easy to skip downstream
// processing (alerting, history recording, etc.) when nothing has changed
// between ticks.
package fingerprint
