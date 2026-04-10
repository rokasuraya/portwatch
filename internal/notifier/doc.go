// Package notifier provides pluggable notification backends for portwatch.
//
// A Notifier can dispatch port-change events to one or more backends
// simultaneously. The following backends are supported:
//
//   - stdout  – always active; writes a human-readable one-liner.
//   - webhook – optional; POSTs a JSON-encoded Event to a configured URL.
//
// Usage:
//
//	n := notifier.New(os.Stdout, "https://hooks.example.com/portwatch")
//	n.Dispatch(notifier.Event{
//	    Timestamp: time.Now(),
//	    Opened:    []string{"tcp:8080"},
//	    Closed:    []string{},
//	})
package notifier
