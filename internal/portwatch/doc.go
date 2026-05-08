// Package portwatch provides the top-level Watcher and Runner types that
// tie together scanning, state management, and alerting into a single
// cohesive monitoring unit.
//
// Typical usage:
//
//	w, err := portwatch.New(cfg, stateFile)
//	if err != nil { ... }
//
//	runner := portwatch.NewRunner(w, cfg.ScanInterval, nil)
//	if err := runner.Run(ctx); err != nil { ... }
//
// The Runner drives the Watcher on a fixed interval and logs transient
// errors without stopping the loop, so the daemon remains resilient to
// momentary scan failures.
package portwatch
