// Package debounce provides a debouncer that delays execution of a callback
// until a quiet period has elapsed since the last trigger call.
//
// This is useful for coalescing rapid bursts of port-change events into a
// single notification, avoiding alert storms when many ports open or close
// in quick succession.
//
// Usage:
//
//	d := debounce.New(500*time.Millisecond, func() {
//		fmt.Println("fired")
//	})
//	d.Trigger() // resets the timer each call
//	d.Flush()   // cancel any pending timer immediately
package debounce
