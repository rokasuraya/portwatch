// Package fanout provides a generic, named-subscriber fan-out broadcaster.
//
// A Fanout distributes published values to all registered subscribers
// concurrently. Each subscriber is identified by a unique name so it can be
// replaced or removed independently.
//
// Example usage:
//
//	f := fanout.New[string]()
//	f.Subscribe("logger", func(v string) { fmt.Println(v) })
//	f.Publish("port 8080 opened")
//	f.Unsubscribe("logger")
package fanout
