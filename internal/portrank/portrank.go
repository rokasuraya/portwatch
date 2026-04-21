// Package portrank ranks open ports by a composite risk score derived from
// classification severity, observed frequency, and port age.
package portrank

import (
	"sort"
	"sync"

	"portwatch/internal/snapshot"
)

// Entry holds a ranked port entry.
type Entry struct {
	Port     int
	Protocol string
	Score    float64
	Reasons  []string
}

// Ranker scores and sorts snapshot entries.
type Ranker struct {
	mu      sync.Mutex
	scorers []ScorerFunc
}

// ScorerFunc assigns a partial score and optional reason to a snapshot entry.
type ScorerFunc func(e snapshot.Entry) (float64, string)

// New returns a Ranker with the provided scorers applied in order.
func New(scorers ...ScorerFunc) *Ranker {
	return &Ranker{scorers: scorers}
}

// Rank evaluates all entries in snap and returns them sorted by descending score.
func (r *Ranker) Rank(snap *snapshot.Snapshot) []Entry {
	if snap == nil {
		return nil
	}

	r.mu.Lock()
	scorers := r.scorers
	r.mu.Unlock()

	entries := snap.Entries()
	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		var total float64
		var reasons []string
		for _, fn := range scorers {
			s, reason := fn(e)
			total += s
			if reason != "" {
				reasons = append(reasons, reason)
			}
		}
		result = append(result, Entry{
			Port:     e.Port,
			Protocol: e.Protocol,
			Score:    total,
			Reasons:  reasons,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result
}

// AddScorer appends a scorer to the ranker at runtime.
func (r *Ranker) AddScorer(fn ScorerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.scorers = append(r.scorers, fn)
}
