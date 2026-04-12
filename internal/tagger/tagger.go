// Package tagger assigns human-readable labels to scanned port entries
// based on well-known port-to-service mappings.
package tagger

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// well-known maps port numbers to common service names.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger labels scanner entries with a service name.
type Tagger struct {
	extra map[int]string
}

// New returns a Tagger. Additional port-to-label mappings may be supplied
// via extra; they take precedence over the built-in well-known list.
func New(extra map[int]string) *Tagger {
	if extra == nil {
		extra = make(map[int]string)
	}
	return &Tagger{extra: extra}
}

// Label returns the service name for the given port, or a generic
// "port/<n>" label when no mapping is found.
func (t *Tagger) Label(port int) string {
	if name, ok := t.extra[port]; ok {
		return name
	}
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

// Tag annotates a slice of scanner entries in-place, setting each
// entry's Label field to the resolved service name.
func (t *Tagger) Tag(entries []scanner.Entry) {
	for i := range entries {
		entries[i].Label = t.Label(entries[i].Port)
	}
}
