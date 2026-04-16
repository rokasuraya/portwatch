// Package portquota provides a threshold guard for the total number of
// concurrently open ports observed in a snapshot.
//
// A Quota is created with a maximum port count. Each call to Check
// compares the snapshot length against that maximum and writes a
// human-readable warning to the configured writer on first breach.
// Subsequent calls suppress the warning until the count drops back
// below the limit, avoiding log spam during sustained overages.
package portquota
