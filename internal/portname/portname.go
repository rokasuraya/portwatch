// Package portname maps port numbers to human-readable service names.
// It provides a fast in-memory lookup with support for custom overrides.
package portname

import "fmt"

// Mapper resolves port numbers to service names.
type Mapper struct {
	builtIn  map[string]string
	custom   map[string]string
}

// well-known port-to-service mappings (protocol:port -> name).
var builtInNames = map[string]string{
	"tcp:22":   "ssh",
	"tcp:23":   "telnet",
	"tcp:25":   "smtp",
	"tcp:53":   "dns",
	"udp:53":   "dns",
	"tcp:80":   "http",
	"tcp:110":  "pop3",
	"tcp:143":  "imap",
	"tcp:443":  "https",
	"tcp:3306": "mysql",
	"tcp:5432": "postgresql",
	"tcp:6379": "redis",
	"tcp:8080": "http-alt",
	"tcp:8443": "https-alt",
	"tcp:27017": "mongodb",
}

// New returns a Mapper loaded with built-in names and any custom overrides.
// A nil or empty overrides map is valid.
func New(overrides map[string]string) *Mapper {
	custom := make(map[string]string, len(overrides))
	for k, v := range overrides {
		custom[k] = v
	}
	return &Mapper{
		builtIn: builtInNames,
		custom:  custom,
	}
}

// Lookup returns the service name for the given protocol and port.
// Custom overrides take precedence over built-in names.
// If no name is found, an empty string is returned.
func (m *Mapper) Lookup(protocol string, port int) string {
	k := key(protocol, port)
	if name, ok := m.custom[k]; ok {
		return name
	}
	return m.builtIn[k]
}

// Register adds or replaces a custom mapping at runtime.
func (m *Mapper) Register(protocol string, port int, name string) {
	m.custom[key(protocol, port)] = name
}

// LookupWithFallback returns the service name or a formatted fallback string
// such as "port/3000" when no mapping exists.
func (m *Mapper) LookupWithFallback(protocol string, port int) string {
	if name := m.Lookup(protocol, port); name != "" {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

func key(protocol string, port int) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
