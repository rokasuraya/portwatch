package history

import (
	"errors"

	"github.com/user/portwatch/internal/state"
)

// Collector bridges a state.Diff and a History, recording every
// opened/closed port as a discrete event.
type Collector struct {
	h *History
}

// NewCollector wraps h in a Collector.
func NewCollector(h *History) *Collector {
	return &Collector{h: h}
}

// Collect records all diffs from d into the underlying History.
// All Record errors are joined and returned so that a single failure
// does not prevent the remaining events from being persisted.
func (c *Collector) Collect(d state.Diff) error {
	var errs []error
	for _, e := range d.Opened {
		if err := c.h.Record("opened", e.Proto, e.Port); err != nil {
			errs = append(errs, err)
		}
	}
	for _, e := range d.Closed {
		if err := c.h.Record("closed", e.Proto, e.Port); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
