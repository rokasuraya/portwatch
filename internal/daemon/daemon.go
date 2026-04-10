package daemon

import (
	"context"
	"log"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

// Daemon orchestrates periodic port scanning and alerting.
type Daemon struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	state   *state.State
	alerter *alert.Alerter
}

// New creates a new Daemon with the provided configuration.
func New(cfg *config.Config, st *state.State, al *alert.Alerter) *Daemon {
	sc := scanner.New(cfg.Timeout)
	return &Daemon{
		cfg:     cfg,
		scanner: sc,
		state:   st,
		alerter: al,
	}
}

// Run starts the daemon loop, scanning ports at the configured interval.
// It blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon starting (interval: %s, ports: %d-%d)",
		d.cfg.Interval, d.cfg.PortStart, d.cfg.PortEnd)

	if err := d.tick(); err != nil {
		log.Printf("initial scan error: %v", err)
	}

	ticker := time.NewTicker(d.cfg.GetScanDuration())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("scan error: %v", err)
			}
		}
	}
}

// tick performs a single scan cycle: scan ports, update state, and alert on changes.
func (d *Daemon) tick() error {
	open, err := d.scanner.ScanRange(d.cfg.PortStart, d.cfg.PortEnd)
	if err != nil {
		return err
	}

	diff, err := d.state.Update(open)
	if err != nil {
		return err
	}

	return d.alerter.Notify(diff)
}
