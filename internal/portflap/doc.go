// Package portflap provides a flap detector for monitored ports.
//
// A port is considered to be "flapping" when it repeatedly opens and closes
// within a short observation window. This is a common symptom of unstable
// services, misconfigured firewalls, or transient network conditions.
//
// Usage:
//
//	d := portflap.New(4, 2*time.Minute)
//	d.Observe(diff.Opened, diff.Closed)
package portflap
