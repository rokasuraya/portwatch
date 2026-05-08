// Package portclaim provides a thread-safe registry that maps port+protocol
// pairs to ownership metadata (process name, PID). It is intended to be
// populated by an OS-level scanner (e.g. reading /proc/net or calling lsof)
// and queried during diff annotation so that alerts can include the name of
// the process that opened or closed a port.
package portclaim
