// Package grace implements a graceful shutdown coordinator for portwatch.
//
// A Coordinator allows components to register in-flight work via Acquire
// and Release, and provides a Shutdown method that blocks until all
// registered work completes or a configurable timeout elapses.
//
// Typical usage:
//
//	g := grace.New(5 * time.Second)
//
//	// Inside a goroutine doing work:
//	if !g.Acquire() {
//		return // shutdown already started
//	}
//	defer g.Release()
//
//	// On SIGTERM / context cancel:
//	if err := g.Shutdown(ctx); err != nil {
//		log.Println("shutdown timed out:", err)
//	}
package grace
