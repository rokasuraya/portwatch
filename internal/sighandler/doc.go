// Package sighandler wraps os/signal to provide a clean, testable interface
// for reacting to OS signals (SIGINT, SIGTERM) in the portwatch daemon.
//
// Typical usage:
//
//	sh := sighandler.New()
//	ctx, cancel := sh.WithCancel(context.Background())
//	defer cancel()
//	// pass ctx to daemon.Run — it will stop cleanly on Ctrl-C
package sighandler
