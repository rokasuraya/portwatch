// Package portevict provides a quiet-period guard for ports that have
// been forcibly removed from the monitored set.
//
// When a port is evicted it enters a configurable quiet period during
// which re-appearance of that port is silently ignored.  Once the quiet
// period expires the port is treated as new again and will trigger the
// normal alert pipeline.
//
// Typical usage:
//
//	evict := portevict.New(2 * time.Minute)
//	evict.Evict(entry)
//	if !evict.IsEvicted(entry) { /* alert */ }
package portevict
