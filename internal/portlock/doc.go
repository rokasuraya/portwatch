// Package portlock allows specific port/protocol pairs to be "locked" so that
// any detected changes on those ports are suppressed from the alerting and
// reporting pipeline.
//
// A locked port is not the same as an acknowledged or baselined port — it is
// an explicit operator decision to permanently silence a port for the lifetime
// of the daemon process.
//
// Usage:
//
//	l := portlock.New(os.Stderr)
//	l.Lock(22, "tcp", "ssh always open")
//	if l.IsLocked(port, proto) {
//	    // skip alerting
//	}
package portlock
