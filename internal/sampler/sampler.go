// Package sampler provides periodic port-scan sampling with jitter to
// avoid thundering-herd effects when multiple portwatch instances run
// on the same host.
package sampler

import (
	"context"
	"math/rand"
	"time"
)

// Scanner is the minimal interface required to collect a port snapshot.
type Scanner interface {
	ScanPortRange(start, end int) ([]string, error)
}

// Sampler wraps a Scanner and fires a callback on every successful sample.
type Sampler struct {
	scanner  Scanner
	interval time.Duration
	jitter   time.Duration
	onSample func(ports []string)
}

// New creates a Sampler.
//
// interval is the base cadence between scans.
// jitter is the maximum random offset added to each interval (0 disables jitter).
// onSample is called with the collected port list after every successful scan.
func New(scanner Scanner, interval, jitter time.Duration, onSample func([]string)) *Sampler {
	return &Sampler{
		scanner:  scanner,
		interval: interval,
		jitter:   jitter,
		onSample: onSample,
	}
}

// Run blocks until ctx is cancelled, sampling at each jittered interval.
func (s *Sampler) Run(ctx context.Context, startPort, endPort int) error {
	for {
		wait := s.interval
		if s.jitter > 0 {
			//nolint:gosec // non-cryptographic jitter is intentional
			wait += time.Duration(rand.Int63n(int64(s.jitter)))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}

		ports, err := s.scanner.ScanPortRange(startPort, endPort)
		if err != nil {
			// Non-fatal: log would go here in a real implementation.
			continue
		}
		if s.onSample != nil {
			s.onSample(ports)
		}
	}
}
