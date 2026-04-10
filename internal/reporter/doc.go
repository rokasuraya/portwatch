// Package reporter provides utilities for formatting and writing portwatch
// scan reports. Reports can be emitted as human-readable text lines or as
// newline-delimited JSON (NDJSON), suitable for log aggregation pipelines.
//
// Usage:
//
//	r, err := reporter.New("/var/log/portwatch.log", true)
//	if err != nil { ... }
//	report := reporter.BuildReport(opened, closed, currentOpen)
//	r.Write(report)
package reporter
