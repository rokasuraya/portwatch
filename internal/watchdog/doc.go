// Package watchdog implements a heartbeat-based stall detector for the
// portwatch scan loop.
//
// Usage:
//
//	wd := watchdog.New(2*time.Second, os.Stderr)
//	go wd.Run(ctx, cfg.ScanInterval)
//
//	// Inside the scan loop:
//	wd.Beat()
//
// If the scan loop fails to call Beat within interval+tolerance the watchdog
// writes a human-readable warning to the configured writer.
package watchdog
