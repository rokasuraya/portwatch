// Package eventbus provides a generic publish-subscribe mechanism for
// broadcasting typed events to multiple named subscribers.
//
// A Bus[T] maintains a set of named handlers. Publishing an event delivers
// it synchronously to every registered handler. Handlers can be added or
// removed at any time; replacing a handler by name is idempotent.
//
// Example usage:
//
//	bus := eventbus.New[scanner.Entry]()
//	bus.Subscribe("logger", func(e scanner.Entry) { log.Println(e) })
//	bus.Publish(entry)
package eventbus
