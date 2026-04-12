package audit

import "github.com/user/portwatch/internal/snapshot"

// Collector feeds snapshot diffs into an Audit log.
type Collector struct {
	audit *Audit
}

// NewCollector returns a Collector backed by the given Audit.
func NewCollector(a *Audit) *Collector {
	return &Collector{audit: a}
}

// Collect records opened and closed port events from a snapshot diff.
func (c *Collector) Collect(opened, closed []snapshot.Entry) {
	for _, e := range opened {
		_ = c.audit.LogOpened(e.Protocol, e.Port)
	}
	for _, e := range closed {
		_ = c.audit.LogClosed(e.Protocol, e.Port)
	}
}
