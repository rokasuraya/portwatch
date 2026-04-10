// Package throttle provides a concurrency-safe rate-limiter for portwatch alerts.
//
// When a port repeatedly transitions between open and closed states (e.g. due
// to a flapping service), the alert subsystem can generate a flood of
// notifications. Throttle suppresses repeated alerts for the same key until a
// configurable cooldown window has elapsed.
//
// Typical usage:
//
//	th := throttle.New(5 * time.Minute)
//	if th.Allow("tcp:8080") {
//	    alerter.Notify(diff)
//	}
//
// Call Purge periodically (e.g. once per daemon tick) to reclaim memory from
// expired entries accumulated during long-running sessions.
package throttle
