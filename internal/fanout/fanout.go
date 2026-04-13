// Package fanout wires an eventbus.Bus to the rest of the portwatch
// pipeline, broadcasting snapshot diffs to all registered consumers.
package fanout

import (
	"context"
	"log"

	"portwatch/internal/eventbus"
	"portwatch/internal/snapshot"
)

// Diff carries the opened and closed port entries produced by a single tick.
type Diff struct {
	Opened []snapshot.Entry
	Closed []snapshot.Entry
}

// Fanout distributes Diff events to registered subscribers via an eventbus.
type Fanout struct {
	bus *eventbus.Bus[Diff]
	log *log.Logger
}

// New returns a Fanout backed by the provided logger.
// If logger is nil, output is discarded.
func New(logger *log.Logger) *Fanout {
	if logger == nil {
		logger = log.New(log.Writer(), "", 0)
	}
	return &Fanout{
		bus: eventbus.New[Diff](),
		log: logger,
	}
}

// Subscribe registers a named handler that will be called for every Diff
// published via Broadcast. Registering the same name twice replaces the
// previous handler.
func (f *Fanout) Subscribe(name string, handler func(Diff)) {
	f.bus.Subscribe(name, handler)
}

// Unsubscribe removes the handler identified by name. It is a no-op if the
// name is not registered.
func (f *Fanout) Unsubscribe(name string) {
	f.bus.Unsubscribe(name)
}

// Broadcast publishes d to all registered subscribers. It returns immediately
// once all handlers have been called. If ctx is already done the broadcast is
// skipped.
func (f *Fanout) Broadcast(ctx context.Context, d Diff) {
	select {
	case <-ctx.Done():
		f.log.Printf("fanout: context done, skipping broadcast")
		return
	default:
	}
	f.bus.Publish(d)
}

// Len returns the number of currently registered subscribers.
func (f *Fanout) Len() int {
	return f.bus.Len()
}
