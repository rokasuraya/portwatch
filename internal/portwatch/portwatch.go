// Package portwatch wires together the core components of the portwatch daemon
// into a single runnable unit.
package portwatch

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"portwatch/internal/config"
	"portwatch/internal/pipeline"
	"portwatch/internal/scanner"
	"portwatch/internal/snapshot"
	"portwatch/internal/ticker"
)

// Watcher runs the port-monitoring loop until the context is cancelled.
type Watcher struct {
	cfg    *config.Config
	store  *snapshot.Store
	pipe   *pipeline.Pipeline
	tick   *ticker.Ticker
	out    io.Writer
}

// New constructs a Watcher from the supplied config.
// stateFile is the path used to persist snapshot state across restarts.
func New(cfg *config.Config, stateFile string) (*Watcher, error) {
	if cfg == nil {
		return nil, fmt.Errorf("portwatch: config must not be nil")
	}

	sc := scanner.New(cfg.Timeout)
	store, err := snapshot.NewStore(stateFile)
	if err != nil {
		return nil, fmt.Errorf("portwatch: open store: %w", err)
	}

	pipe := pipeline.New(sc, store)

	interval := cfg.GetScanDuration()
	tk := ticker.New(interval, 0)

	return &Watcher{
		cfg:   cfg,
		store: store,
		pipe:  pipe,
		tick:  tk,
		out:   os.Stdout,
	}, nil
}

// SetOutput redirects human-readable status lines (default: os.Stdout).
func (w *Watcher) SetOutput(out io.Writer) {
	w.out = out
}

// Run starts the ticker-driven scan loop and blocks until ctx is done.
func (w *Watcher) Run(ctx context.Context) error {
	fmt.Fprintf(w.out, "portwatch: starting — interval %s, ports %d-%d\n",
		w.cfg.GetScanDuration(), w.cfg.StartPort, w.cfg.EndPort)

	return w.tick.Run(ctx, func(at time.Time) {
		d, err := w.pipe.Tick(ctx, w.cfg.StartPort, w.cfg.EndPort)
		if err != nil {
			fmt.Fprintf(w.out, "portwatch: scan error at %s: %v\n", at.Format(time.RFC3339), err)
			return
		}
		fmt.Fprintf(w.out, "portwatch: scan complete at %s in %s\n",
			at.Format(time.RFC3339), d)
	})
}
