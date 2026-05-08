// Package portlabel assigns human-readable labels to port+protocol pairs
// based on a built-in registry and optional user-supplied overrides.
package portlabel

import "fmt"

// Label holds the display name and category for a port.
type Label struct {
	Name     string
	Category string
}

// Labeler maps port+protocol pairs to Labels.
type Labeler struct {
	entries map[string]Label
}

var builtIn = map[string]Label{
	"22/tcp":   {Name: "SSH", Category: "remote-access"},
	"23/tcp":   {Name: "Telnet", Category: "remote-access"},
	"25/tcp":   {Name: "SMTP", Category: "mail"},
	"53/tcp":   {Name: "DNS", Category: "infrastructure"},
	"53/udp":   {Name: "DNS", Category: "infrastructure"},
	"80/tcp":   {Name: "HTTP", Category: "web"},
	"443/tcp":  {Name: "HTTPS", Category: "web"},
	"3306/tcp": {Name: "MySQL", Category: "database"},
	"5432/tcp": {Name: "PostgreSQL", Category: "database"},
	"6379/tcp": {Name: "Redis", Category: "cache"},
	"8080/tcp": {Name: "HTTP-Alt", Category: "web"},
	"3389/tcp": {Name: "RDP", Category: "remote-access"},
}

// New returns a Labeler seeded with built-in labels and any overrides.
// A nil overrides map is safe.
func New(overrides map[string]Label) *Labeler {
	l := &Labeler{entries: make(map[string]Label, len(builtIn)+len(overrides))}
	for k, v := range builtIn {
		l.entries[k] = v
	}
	for k, v := range overrides {
		l.entries[k] = v
	}
	return l
}

// Lookup returns the Label for the given port and protocol.
// If no entry exists, a generic label is returned with ok=false.
func (l *Labeler) Lookup(port uint16, protocol string) (Label, bool) {
	k := key(port, protocol)
	if lbl, ok := l.entries[k]; ok {
		return lbl, true
	}
	return Label{Name: fmt.Sprintf("port-%d", port), Category: "unknown"}, false
}

// Register adds or replaces a label entry at runtime.
func (l *Labeler) Register(port uint16, protocol string, lbl Label) {
	l.entries[key(port, protocol)] = lbl
}

// Len returns the number of registered entries.
func (l *Labeler) Len() int { return len(l.entries) }

func key(port uint16, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}
