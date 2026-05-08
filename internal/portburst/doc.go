// Package portburst provides burst detection for port-open events.
//
// A Detector maintains a rolling time window of port-open timestamps.
// When the number of openings within that window reaches or exceeds the
// configured threshold a warning is written to the configured writer.
//
// Typical usage:
//
//	det := portburst.New(10, 30*time.Second, os.Stderr)
//	det.Observe(diff)  // called after each scan tick
package portburst
