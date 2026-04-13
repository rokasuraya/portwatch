// Package labelmap provides a registry that maps port numbers to
// human-readable service names, merging a built-in table with any
// caller-supplied overrides.
package labelmap

import "fmt"

// builtIn contains a small subset of well-known IANA service names.
var builtIn = map[uint16]string{
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

// LabelMap holds the merged port-to-label mapping.
type LabelMap struct {
	table map[uint16]string
}

// New returns a LabelMap that merges builtIn with extra.
// Values in extra take precedence over builtIn.
func New(extra map[uint16]string) *LabelMap {
	table := make(map[uint16]string, len(builtIn)+len(extra))
	for k, v := range builtIn {
		table[k] = v
	}
	for k, v := range extra {
		table[k] = v
	}
	return &LabelMap{table: table}
}

// Lookup returns the service label for port, and whether it was found.
func (lm *LabelMap) Lookup(port uint16) (string, bool) {
	v, ok := lm.table[port]
	return v, ok
}

// Label returns the service label for port, or a generic "port/<n>" string
// when no label is registered.
func (lm *LabelMap) Label(port uint16) string {
	if v, ok := lm.table[port]; ok {
		return v
	}
	return fmt.Sprintf("port/%d", port)
}

// Register adds or replaces the label for port at runtime.
func (lm *LabelMap) Register(port uint16, label string) {
	lm.table[port] = label
}

// Len returns the number of entries in the map.
func (lm *LabelMap) Len() int {
	return len(lm.table)
}
