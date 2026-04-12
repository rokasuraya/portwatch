// Package resolver provides port-to-service name resolution for portwatch.
//
// It maintains a built-in table of well-known port/protocol pairs and allows
// callers to supply additional mappings at construction time or register them
// dynamically via Register. Lookups are safe for concurrent use.
package resolver
