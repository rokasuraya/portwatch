// Package portcap enforces a maximum number of simultaneously open ports.
//
// A PortCap is configured with a maximum port count. On each call to Check,
// it compares the number of entries in the provided snapshot against the
// configured limit and writes a warning for every entry that exceeds it.
//
// Usage:
//
//	c := portcap.New(50, os.Stderr)
//	violations := c.Check(snap)
//
// Runner integrates PortCap into a periodic loop, calling a SnapshotFunc on
// every tick and forwarding results to Check.
package portcap
