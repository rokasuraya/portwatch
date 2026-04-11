// Package sighandler provides OS signal handling for graceful shutdown.
package sighandler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and cancels a context on receipt.
type Handler struct {
	signals []os.Signal
	notify  func(chan<- os.Signal, ...os.Signal)
	stop    func(chan<- os.Signal)
}

// New returns a Handler that reacts to the provided signals.
// If no signals are given, SIGINT and SIGTERM are used by default.
func New(sigs ...os.Signal) *Handler {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return &Handler{
		signals: sigs,
		notify:  signal.Notify,
		stop:    signal.Stop,
	}
}

// Wait blocks until one of the registered signals is received or the parent
// context is cancelled. It returns the signal that triggered the shutdown, or
// nil if the parent context expired first.
func (h *Handler) Wait(ctx context.Context) os.Signal {
	ch := make(chan os.Signal, 1)
	h.notify(ch, h.signals...)
	defer h.stop(ch)

	select {
	case sig := <-ch:
		return sig
	case <-ctx.Done():
		return nil
	}
}

// WithCancel returns a derived context that is cancelled when one of the
// registered signals is received. The caller must invoke the returned cancel
// function to release resources.
func (h *Handler) WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		if sig := h.Wait(parent); sig != nil {
			cancel()
		}
	}()
	return ctx, cancel
}
