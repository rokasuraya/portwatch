// Package classify provides risk-level classification for open ports discovered
// by the portwatch scanner.
//
// Entries are assigned one of three levels — low, medium, or high — based on
// well-known port heuristics. Callers may supply a []Rule slice to override the
// built-in defaults for any port/protocol combination.
//
// Usage:
//
//	c := classify.New(nil)          // built-in heuristics only
//	lvl := c.Classify(entry)        // classify a single entry
//	levels := c.Apply(entries)      // classify a whole slice
package classify
