// Package resolver maps port numbers to well-known service names.
package resolver

import (
	"fmt"
	"sync"
)

// Resolver looks up service names for port/protocol pairs.
type Resolver struct {
	mu    sync.RWMutex
	extra map[string]string
}

// builtIn holds a curated set of well-known port/protocol → service mappings.
var builtIn = map[string]string{
	"tcp/21":   "ftp",
	"tcp/22":   "ssh",
	"tcp/23":   "telnet",
	"tcp/25":   "smtp",
	"tcp/53":   "dns",
	"udp/53":   "dns",
	"tcp/80":   "http",
	"tcp/110":  "pop3",
	"tcp/143":  "imap",
	"tcp/443":  "https",
	"tcp/3306": "mysql",
	"tcp/5432": "postgresql",
	"tcp/6379": "redis",
	"tcp/8080": "http-alt",
	"tcp/8443": "https-alt",
	"tcp/27017": "mongodb",
}

// New returns a Resolver. extra overrides or extends the built-in table.
func New(extra map[string]string) *Resolver {
	r := &Resolver{
		extra: make(map[string]string, len(extra)),
	}
	for k, v := range extra {
		r.extra[k] = v
	}
	return r
}

// Lookup returns the service name for the given port and protocol.
// It returns an empty string when the port is unknown.
func (r *Resolver) Lookup(port uint16, proto string) string {
	k := key(port, proto)
	r.mu.RLock()
	if v, ok := r.extra[k]; ok {
		r.mu.RUnlock()
		return v
	}
	r.mu.RUnlock()
	return builtIn[k]
}

// Register adds or replaces a mapping at runtime.
func (r *Resolver) Register(port uint16, proto, service string) {
	r.mu.Lock()
	r.extra[key(port, proto)] = service
	r.mu.Unlock()
}

func key(port uint16, proto string) string {
	return fmt.Sprintf("%s/%d", proto, port)
}
