// Package classify categorises scanner entries by risk level based on
// port number, protocol, and an optional user-supplied rule set.
package classify

import "github.com/user/portwatch/internal/scanner"

// Level represents the risk classification of an open port.
type Level int

const (
	LevelLow Level = iota
	LevelMedium
	LevelHigh
)

// String returns a human-readable label for the Level.
func (l Level) String() string {
	switch l {
	case LevelHigh:
		return "high"
	case LevelMedium:
		return "medium"
	default:
		return "low"
	}
}

// Rule overrides the default classification for a specific port/protocol pair.
type Rule struct {
	Port     int
	Protocol string
	Level    Level
}

// Classifier assigns risk levels to scanner entries.
type Classifier struct {
	rules map[string]Level
}

// New returns a Classifier loaded with the provided override rules.
// Built-in heuristics apply when no matching rule exists.
func New(rules []Rule) *Classifier {
	c := &Classifier{rules: make(map[string]Level, len(rules))}
	for _, r := range rules {
		c.rules[key(r.Port, r.Protocol)] = r.Level
	}
	return c
}

// Classify returns the risk Level for a single scanner entry.
func (c *Classifier) Classify(e scanner.Entry) Level {
	if lvl, ok := c.rules[key(e.Port, e.Protocol)]; ok {
		return lvl
	}
	return defaultLevel(e.Port)
}

// Apply returns a map of entry → Level for every entry in the slice.
func (c *Classifier) Apply(entries []scanner.Entry) map[scanner.Entry]Level {
	out := make(map[scanner.Entry]Level, len(entries))
	for _, e := range entries {
		out[e] = c.Classify(e)
	}
	return out
}

func key(port int, proto string) string {
	return proto + ":" + itoa(port)
}

func defaultLevel(port int) Level {
	switch {
	case port == 22 || port == 23 || port == 3389:
		return LevelHigh
	case port < 1024:
		return LevelMedium
	default:
		return LevelLow
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
