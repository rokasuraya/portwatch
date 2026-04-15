// Package cooldown provides a thread-safe per-key cooldown tracker.
//
// A Cooldown suppresses repeated events for the same key within a
// configurable quiet period. This is useful for rate-limiting alerts
// or notifications so that a single flapping port does not generate
// a flood of messages.
//
// Example usage:
//
//	cd := cooldown.New(5 * time.Minute)
//	if cd.Allow("port:22:tcp") {
//		// send alert
//	}
package cooldown
