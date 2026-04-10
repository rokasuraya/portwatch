// Package filter implements suppression rules for portwatch alerts.
//
// A Filter is constructed from a slice of Rule values, each specifying a
// port number and protocol ("tcp" or "udp") that should be considered
// expected and therefore excluded from change notifications.
//
// Rules are typically loaded from the portwatch configuration file under
// the "ignore_ports" key and passed to filter.New before the daemon
// processes scanner output.
package filter
