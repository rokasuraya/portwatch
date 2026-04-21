// Package portfreq provides a lightweight frequency tracker for open ports.
//
// Each call to Observe increments an internal counter for every port present
// in the supplied snapshot. The Top method returns the most-frequently-seen
// ports across all observed snapshots, making it easy to distinguish
// persistently open ports from transient noise.
//
// Usage:
//
//	tracker := portfreq.New()
//	tracker.Observe(snap)
//	top5 := tracker.Top(5)
package portfreq
