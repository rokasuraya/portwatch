// Package envelope defines the Envelope type, which bundles a port-scan diff
// with metadata (ID, timestamp, scan duration, and arbitrary labels) for
// consistent handoff between pipeline stages and downstream consumers such as
// notifiers, auditors, and reporters.
package envelope
