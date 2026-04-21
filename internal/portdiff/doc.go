// Package portdiff computes and formats the difference between two port
// snapshots, producing human-readable or JSON output.
//
// Basic usage:
//
//	diff := portdiff.Compute(prevSnap, nextSnap, labeler)
//	if !diff.IsEmpty() {
//		portdiff.Format(os.Stdout, diff)
//	}
//
// The Labeler interface can be satisfied by any type that maps a port number
// and protocol string to a descriptive name (e.g. internal/labelmap).
package portdiff
