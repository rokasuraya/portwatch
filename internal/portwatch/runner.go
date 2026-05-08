package portwatch

import (
	"context"
	"log"
	"time"
)

// Runner periodically invokes the Watcher and logs any errors.
type Runner struct {
	watcher  *Watcher
	interval time.Duration
	log      *log.Logger
}

// NewRunner creates a Runner that ticks the watcher at the given interval.
// If logger is nil it defaults to the standard logger.
func NewRunner(w *Watcher, interval time.Duration, logger *log.Logger) *Runner {
	if logger == nil {
		logger = log.Default()
	}
	return &Runner{
		watcher:  w,
		interval: interval,
		log:      logger,
	}
}

// Run blocks until ctx is cancelled, invoking Tick on each interval.
func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.watcher.Tick(ctx); err != nil {
				r.log.Printf("portwatch: tick error: %v", err)
			}
		}
	}
}
