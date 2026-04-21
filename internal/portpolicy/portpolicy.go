// Package portpolicy enforces allow/deny policies on observed ports.
// A policy is a named rule that either permits or rejects a port+protocol
// combination. Rules are evaluated in insertion order; the first match wins.
// If no rule matches the default action is Allow.
package portpolicy

import (
	"fmt"
	"sync"

	"portwatch/internal/snapshot"
)

// Action describes whether a port is permitted or denied.
type Action int

const (
	Allow Action = iota
	Deny
)

func (a Action) String() string {
	if a == Deny {
		return "deny"
	}
	return "allow"
}

// Rule is a single policy entry.
type Rule struct {
	Name     string
	Port     int
	Protocol string // "tcp" | "udp" | "" (matches any)
	Action   Action
}

// Violation is produced when a Deny rule matches an entry.
type Violation struct {
	Rule  string
	Entry snapshot.Entry
}

func (v Violation) Error() string {
	return fmt.Sprintf("policy %q denies %s/%d", v.Rule, v.Entry.Protocol, v.Entry.Port)
}

// Policy holds an ordered list of rules.
type Policy struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns an empty Policy.
func New() *Policy { return &Policy{} }

// Add appends a rule to the policy.
func (p *Policy) Add(r Rule) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, r)
}

// Evaluate checks a single entry against all rules.
// It returns (Allow, nil) when no deny rule matches.
func (p *Policy) Evaluate(e snapshot.Entry) (Action, *Violation) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.rules {
		if r.Port != e.Port {
			continue
		}
		if r.Protocol != "" && r.Protocol != e.Protocol {
			continue
		}
		if r.Action == Deny {
			v := &Violation{Rule: r.Name, Entry: e}
			return Deny, v
		}
		return Allow, nil
	}
	return Allow, nil
}

// Check evaluates every entry in a snapshot and returns all violations.
func (p *Policy) Check(snap *snapshot.Snapshot) []Violation {
	if snap == nil {
		return nil
	}
	var out []Violation
	for _, e := range snap.Entries {
		if _, v := p.Evaluate(e); v != nil {
			out = append(out, *v)
		}
	}
	return out
}
