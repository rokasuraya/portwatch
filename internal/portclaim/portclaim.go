// Package portclaim tracks which process or service has claimed a given port,
// allowing portwatch to annotate diffs with ownership information.
package portclaim

import (
	"fmt"
	"sync"
)

// Claim holds ownership metadata for a single port+protocol pair.
type Claim struct {
	Port     int
	Protocol string
	Owner    string // e.g. process name or service label
	PID      int    // 0 means unknown
}

// Registry stores and retrieves port claims.
type Registry struct {
	mu     sync.RWMutex
	claims map[string]Claim
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		claims: make(map[string]Claim),
	}
}

func claimKey(port int, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}

// Register records a claim for the given port and protocol.
// Calling Register again for the same key overwrites the previous claim.
func (r *Registry) Register(c Claim) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.claims[claimKey(c.Port, c.Protocol)] = c
}

// Lookup returns the Claim for the given port and protocol, and whether it exists.
func (r *Registry) Lookup(port int, protocol string) (Claim, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.claims[claimKey(port, protocol)]
	return c, ok
}

// Release removes the claim for the given port and protocol.
func (r *Registry) Release(port int, protocol string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.claims, claimKey(port, protocol))
}

// Len returns the number of active claims.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.claims)
}
