// Package snapshot provides types and utilities for capturing and comparing
// the set of open ports observed by portwatch at a point in time.
//
// # Overview
//
// A [Snapshot] is an immutable, timestamped list of [Entry] values, each
// representing a single open port. Use [New] to create a snapshot from a
// slice of entries returned by the scanner.
//
// Use [Compare] to produce a [Diff] between two snapshots, identifying
// which ports were opened or closed between observations.
//
// A [Store] wraps a file-backed snapshot so that the daemon can persist
// state across restarts and always have access to both the current and
// previous snapshot without extra bookkeeping.
package snapshot
