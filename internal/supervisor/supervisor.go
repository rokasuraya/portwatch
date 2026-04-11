// Package supervisor coordinates the daemon lifecycle, wiring together
// the scanner, state, alerter, metrics, and healthcheck components so
// the daemon loop remains thin.
package supervisor

import (
	"context"
	"log"
	"time"

	"portwatch/internal/config"
	"portwatch/internal/metrics"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

// Components groups the external dependencies the Supervisor needs.
type Components struct {
	Scanner  *scanner.Scanner
	State    *state.State
	Metrics  *metrics.Metrics
	OnChange func(opened, closed []string)
}

// Supervisor orchestrates a single scan tick and records outcomes.
type Supervisor struct {
	cfg        *config.Config
	components Components
	log        *log.Logger
}

// New returns a Supervisor wired with the provided config and components.
func New(cfg *config.Config, c Components, logger *log.Logger) *Supervisor {
	if logger == nil {
		logger = log.Default()
	}
	return &Supervisor{cfg: cfg, components: c, log: logger}
}

// Tick performs one full scan cycle: scan ports, diff against state,
// invoke the onChange callback for any diff, then persist state.
func (s *Supervisor) Tick(ctx context.Context) error {
	start := time.Now()

	ports, err := s.components.Scanner.ScanPortRange(
		ctx,
		s.cfg.Host,
		s.cfg.PortRangeStart,
		s.cfg.PortRangeEnd,
	)
	if err != nil {
		return err
	}

	opened, closed := s.components.State.Update(ports)
	duration := time.Since(start)

	s.components.Metrics.RecordScan(duration, opened, closed)

	if len(opened) > 0 || len(closed) > 0 {
		if s.components.OnChange != nil {
			s.components.OnChange(opened, closed)
		}
		s.log.Printf("supervisor: opened=%d closed=%d elapsed=%s",
			len(opened), len(closed), duration.Round(time.Millisecond))
	}

	return nil
}
