// Package rollup provides a time-windowed event aggregator for port changes.
//
// During periods of rapid churn — for example when a service restarts and
// briefly closes then reopens a port — rollup collects all opened/closed
// entries within a configurable quiet window and emits a single [Summary]
// instead of one alert per individual change.
//
// Usage:
//
//	r := rollup.New(500*time.Millisecond, func(s rollup.Summary) {
//	    fmt.Printf("opened: %v  closed: %v\n", s.Opened, s.Closed)
//	})
//	r.Add(opened, closed) // call from your scan loop
//
// A zero wait duration disables batching and flushes immediately on every Add.
package rollup
