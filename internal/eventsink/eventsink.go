// Package eventsink provides a fan-out sink that delivers port-change
// events to multiple named handlers in a fire-and-forget fashion.
package eventsink

import (
	"sync"

	"github.com/user/portwatch/internal/snapshot"
)

// Handler is a function that receives a diff produced by a scan tick.
type Handler func(opened, closed []snapshot.Entry)

// Sink delivers scan diffs to a set of registered handlers.
type Sink struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New returns an empty Sink ready for use.
func New() *Sink {
	return &Sink{
		handlers: make(map[string]Handler),
	}
}

// Register adds or replaces a named handler. Passing a nil handler removes
// the entry (equivalent to calling Unregister).
func (s *Sink) Register(name string, h Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if h == nil {
		delete(s.handlers, name)
		return
	}
	s.handlers[name] = h
}

// Unregister removes the handler with the given name. It is a no-op if the
// name is not registered.
func (s *Sink) Unregister(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.handlers, name)
}

// Len returns the number of currently registered handlers.
func (s *Sink) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.handlers)
}

// Emit delivers the diff to every registered handler. Each handler is invoked
// synchronously in an unspecified order. Emit is safe to call concurrently.
func (s *Sink) Emit(opened, closed []snapshot.Entry) {
	s.mu.RLock()
	handlers := make([]Handler, 0, len(s.handlers))
	for _, h := range s.handlers {
		handlers = append(handlers, h)
	}
	s.mu.RUnlock()

	for _, h := range handlers {
		h(opened, closed)
	}
}
