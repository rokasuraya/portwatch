// Package portmute provides temporary muting of port change alerts.
//
// A Muter suppresses notifications for specific port/protocol pairs for a
// configurable duration. Mutes expire automatically; no manual cleanup is
// required. This is useful during planned maintenance windows where known
// port changes should not generate noise.
//
// Example:
//
//	m := portmute.New()
//	m.Mute(8080, "tcp", 30*time.Minute, "maintenance window")
//	if !m.IsMuted(8080, "tcp") {
//		// send alert
//	}
package portmute
