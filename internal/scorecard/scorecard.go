// Package scorecard computes a numeric risk score for a snapshot diff,
// combining severity classification, port count, and trend direction.
package scorecard

import (
	"sync"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/snapshot"
)

// Score holds the computed risk assessment for a single diff.
type Score struct {
	Total    int    // aggregate risk score
	Opened   int    // number of newly opened ports
	Closed   int    // number of newly closed ports
	HighRisk int    // count of high-severity opened ports
	Label    string // human-readable risk band
}

// Scorecard evaluates diffs and produces a Score.
type Scorecard struct {
	mu         sync.Mutex
	classifier *classify.Classifier
	last       Score
}

// New returns a Scorecard backed by the provided Classifier.
func New(c *classify.Classifier) *Scorecard {
	return &Scorecard{classifier: c}
}

// Evaluate computes a Score from opened and closed snapshot entries.
func (s *Scorecard) Evaluate(opened, closed []snapshot.Entry) Score {
	s.mu.Lock()
	defer s.mu.Unlock()

	high := 0
	for _, e := range opened {
		if s.classifier.Classify(e) == classify.High {
			high++
		}
	}

	total := len(opened)*10 + high*15 + len(closed)*2
	label := band(total)

	sc := Score{
		Total:    total,
		Opened:   len(opened),
		Closed:   len(closed),
		HighRisk: high,
		Label:    label,
	}
	s.last = sc
	return sc
}

// Last returns the most recently computed Score.
func (s *Scorecard) Last() Score {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

func band(total int) string {
	switch {
	case total == 0:
		return "none"
	case total < 20:
		return "low"
	case total < 50:
		return "medium"
	default:
		return "high"
	}
}
