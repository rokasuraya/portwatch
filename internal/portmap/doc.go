// Package portmap provides a thread-safe, in-memory registry that tracks
// the current open/closed status of every scanned port.
//
// A PortMap is populated by an Updater which consumes snapshot diffs
// produced by the pipeline and reflects them as live state. Consumers can
// query individual ports via Get, iterate all entries via All, or watch the
// count via Len.
//
// Typical usage:
//
//	pm := portmap.New()
//	updater := portmap.NewUpdater(pm, resolver.Lookup)
//	// inside tick callback:
//	updater.Apply(diff.Opened, diff.Closed)
package portmap
