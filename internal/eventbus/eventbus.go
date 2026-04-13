// Package eventbus provides a simple publish/subscribe mechanism for
// broadcasting port change events to multiple internal consumers.
package eventbus

import "sync"

// Handler is a function that receives a published payload.
type Handler[T any] func(T)

// Bus is a generic, goroutine-safe publish/subscribe event bus.
type Bus[T any] struct {
	mu       sync.RWMutex
	handlers map[string]Handler[T]
}

// New returns an initialised Bus ready for use.
func New[T any]() *Bus[T] {
	return &Bus[T]{
		handlers: make(map[string]Handler[T]),
	}
}

// Subscribe registers a named handler. Registering with the same name
// replaces the previous handler.
func (b *Bus[T]) Subscribe(name string, h Handler[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = h
}

// Unsubscribe removes the handler registered under name. It is a no-op
// if the name is not registered.
func (b *Bus[T]) Unsubscribe(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, name)
}

// Publish delivers payload to every registered handler. Each handler is
// called synchronously in an unspecified order.
func (b *Bus[T]) Publish(payload T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers {
		h(payload)
	}
}

// Len returns the number of currently registered handlers.
func (b *Bus[T]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers)
}
