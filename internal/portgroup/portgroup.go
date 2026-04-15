// Package portgroup provides named groupings of ports for use in rules and reports.
package portgroup

import "fmt"

// Group is a named set of port+protocol pairs.
type Group struct {
	Name    string
	entries map[string]struct{}
}

// Registry holds named port groups.
type Registry struct {
	groups map[string]*Group
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string]*Group)}
}

// Define creates or replaces a named group with the given port/protocol pairs.
// Each entry must be in "port/proto" form, e.g. "22/tcp".
func (r *Registry) Define(name string, entries []string) error {
	g := &Group{Name: name, entries: make(map[string]struct{}, len(entries))}
	for _, e := range entries {
		if e == "" {
			return fmt.Errorf("portgroup: empty entry in group %q", name)
		}
		g.entries[e] = struct{}{}
	}
	r.groups[name] = g
	return nil
}

// Contains reports whether the named group contains the given "port/proto" key.
// Returns false if the group does not exist.
func (r *Registry) Contains(name, key string) bool {
	g, ok := r.groups[name]
	if !ok {
		return false
	}
	_, found := g.entries[key]
	return found
}

// Names returns all defined group names in unspecified order.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.groups))
	for k := range r.groups {
		out = append(out, k)
	}
	return out
}

// Members returns the entries of the named group, or nil if not found.
func (r *Registry) Members(name string) []string {
	g, ok := r.groups[name]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(g.entries))
	for k := range g.entries {
		out = append(out, k)
	}
	return out
}

// Remove deletes the named group. It is a no-op if the group does not exist.
func (r *Registry) Remove(name string) {
	delete(r.groups, name)
}
