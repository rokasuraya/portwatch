package scorecard_test

import (
	"testing"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/scorecard"
	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port uint16, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Proto: proto}
}

func newClassifier() *classify.Classifier {
	return classify.New(nil)
}

func TestNew_ReturnsScorecard(t *testing.T) {
	sc := scorecard.New(newClassifier())
	if sc == nil {
		t.Fatal("expected non-nil Scorecard")
	}
}

func TestEvaluate_EmptyDiff(t *testing.T) {
	sc := scorecard.New(newClassifier())
	s := sc.Evaluate(nil, nil)
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
	if s.Label != "none" {
		t.Errorf("expected label 'none', got %q", s.Label)
	}
}

func TestEvaluate_OpenedPortsIncreasesScore(t *testing.T) {
	sc := scorecard.New(newClassifier())
	opened := []snapshot.Entry{makeEntry(8080, "tcp"), makeEntry(9090, "tcp")}
	s := sc.Evaluate(opened, nil)
	if s.Opened != 2 {
		t.Errorf("expected Opened=2, got %d", s.Opened)
	}
	if s.Total < 20 {
		t.Errorf("expected total >= 20, got %d", s.Total)
	}
}

func TestEvaluate_HighRiskPortBoostsScore(t *testing.T) {
	sc := scorecard.New(newClassifier())
	// port 22 (SSH) is classified High by the default classifier
	opened := []snapshot.Entry{makeEntry(22, "tcp")}
	s := sc.Evaluate(opened, nil)
	if s.HighRisk != 1 {
		t.Errorf("expected HighRisk=1, got %d", s.HighRisk)
	}
	// 10 (opened) + 15 (high) = 25 → medium
	if s.Label != "medium" {
		t.Errorf("expected label 'medium', got %q", s.Label)
	}
}

func TestEvaluate_ClosedPortsAddSmallPenalty(t *testing.T) {
	sc := scorecard.New(newClassifier())
	closed := []snapshot.Entry{makeEntry(80, "tcp")}
	s := sc.Evaluate(nil, closed)
	if s.Closed != 1 {
		t.Errorf("expected Closed=1, got %d", s.Closed)
	}
	if s.Total != 2 {
		t.Errorf("expected total 2, got %d", s.Total)
	}
	if s.Label != "low" {
		t.Errorf("expected label 'low', got %q", s.Label)
	}
}

func TestLast_ReturnsLatestScore(t *testing.T) {
	sc := scorecard.New(newClassifier())
	sc.Evaluate([]snapshot.Entry{makeEntry(443, "tcp")}, nil)
	last := sc.Last()
	if last.Opened != 1 {
		t.Errorf("expected Last().Opened=1, got %d", last.Opened)
	}
}

func TestEvaluate_HighLabelThreshold(t *testing.T) {
	sc := scorecard.New(newClassifier())
	// 4 high-risk ports: 4*10 + 4*15 = 100 → high
	opened := []snapshot.Entry{
		makeEntry(22, "tcp"),
		makeEntry(3389, "tcp"),
		makeEntry(23, "tcp"),
		makeEntry(445, "tcp"),
	}
	s := sc.Evaluate(opened, nil)
	if s.Label != "high" {
		t.Errorf("expected label 'high', got %q", s.Label)
	}
}
