// Package pipeline wires together the scan-to-alert flow as a reusable
// processing stage that can be composed by higher-level runners.
package pipeline

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// ScanFunc performs a port scan and returns the resulting entries.
type ScanFunc func(ctx context.Context) ([]snapshot.Entry, error)

// StageFunc processes a diff produced by comparing two snapshots.
type StageFunc func(ctx context.Context, opened, closed []snapshot.Entry) error

// Pipeline orchestrates a single scan tick: scan → compare → notify stages.
type Pipeline struct {
	scan   ScanFunc
	store  *snapshot.Store
	stages []StageFunc
}

// New returns a Pipeline wired with the provided scan function, snapshot
// store, and zero or more processing stages.
func New(scan ScanFunc, store *snapshot.Store, stages ...StageFunc) *Pipeline {
	return &Pipeline{
		scan:   scan,
		store:  store,
		stages: stages,
	}
}

// Tick executes one full scan cycle. It scans, persists the snapshot, computes
// the diff, and fans out to every registered stage. The first stage error is
// returned; subsequent stages are still executed.
func (p *Pipeline) Tick(ctx context.Context) (time.Duration, error) {
	start := time.Now()

	entries, err := p.scan(ctx)
	if err != nil {
		return 0, err
	}

	next := snapshot.New(entries)
	opened, closed := p.diff(next)

	if err := p.store.Set(next); err != nil {
		return 0, err
	}

	var firstErr error
	for _, stage := range p.stages {
		if sErr := stage(ctx, opened, closed); sErr != nil && firstErr == nil {
			firstErr = sErr
		}
	}

	return time.Since(start), firstErr
}

// diff returns the opened/closed entry slices relative to the current stored
// snapshot. If no previous snapshot exists both slices are empty.
func (p *Pipeline) diff(next *snapshot.Snapshot) (opened, closed []snapshot.Entry) {
	prev := p.store.Current()
	if prev == nil {
		return nil, nil
	}
	return snapshot.Compare(prev, next)
}
