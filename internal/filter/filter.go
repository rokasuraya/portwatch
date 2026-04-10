// Package filter provides port filtering logic for portwatch.
// It allows users to define rules that suppress alerts for known/expected ports.
package filter

import "fmt"

// Rule represents a single filter rule.
type Rule struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"` // "tcp" or "udp"
	Comment  string `json:"comment,omitempty"`
}

// Filter holds a set of rules and applies them to port lists.
type Filter struct {
	rules map[string]struct{}
}

// New creates a Filter from the provided rules.
func New(rules []Rule) *Filter {
	f := &Filter{
		rules: make(map[string]struct{}, len(rules)),
	}
	for _, r := range rules {
		f.rules[key(r.Port, r.Protocol)] = struct{}{}
	}
	return f
}

// Allow returns true if the port/protocol combination is NOT suppressed.
func (f *Filter) Allow(port int, protocol string) bool {
	_, suppressed := f.rules[key(port, protocol)]
	return !suppressed
}

// Apply filters a slice of "port/protocol" strings, returning only allowed entries.
// Each entry is expected in the format returned by scanner (e.g. "8080/tcp").
func (f *Filter) Apply(entries []string) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		var port int
		var proto string
		if _, err := fmt.Sscanf(e, "%d/%s", &port, &proto); err != nil {
			out = append(out, e) // keep unparseable entries
			continue
		}
		if f.Allow(port, proto) {
			out = append(out, e)
		}
	}
	return out
}

func key(port int, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}
