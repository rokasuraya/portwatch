// Package history provides a persistent, capped event log of port
// changes observed by portwatch.
//
// Usage:
//
//	h, err := history.New("/var/lib/portwatch/history.json", 500)
//	if err != nil { ... }
//
//	collector := history.NewCollector(h)
//	// After each scan tick:
//	if err := collector.Collect(diff); err != nil { ... }
//
// Events are stored as JSON and survive process restarts. The log is
// automatically trimmed to the configured capacity so disk usage
// remains bounded.
package history
