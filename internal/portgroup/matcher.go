package portgroup

import "fmt"

// Matcher checks whether a port+protocol pair belongs to any registered group.
type Matcher struct {
	registry *Registry
}

// NewMatcher returns a Matcher backed by the given Registry.
func NewMatcher(r *Registry) *Matcher {
	return &Matcher{registry: r}
}

// MatchResult holds the groups that matched a given key.
type MatchResult struct {
	Key    string
	Groups []string
}

// Match returns all group names that contain the given port and protocol.
// port and proto are combined as "port/proto", e.g. "443/tcp".
func (m *Matcher) Match(port uint16, proto string) MatchResult {
	key := fmt.Sprintf("%d/%s", port, proto)
	var matched []string
	for _, name := range m.registry.Names() {
		if m.registry.Contains(name, key) {
			matched = append(matched, name)
		}
	}
	return MatchResult{Key: key, Groups: matched}
}

// AnyMatch returns true if the port+protocol belongs to at least one group.
func (m *Matcher) AnyMatch(port uint16, proto string) bool {
	return len(m.Match(port, proto).Groups) > 0
}

// InGroup returns true if the port+protocol is a member of the named group.
func (m *Matcher) InGroup(name string, port uint16, proto string) bool {
	key := fmt.Sprintf("%d/%s", port, proto)
	return m.registry.Contains(name, key)
}
