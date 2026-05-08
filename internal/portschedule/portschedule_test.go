package portschedule_test

import (
	"testing"
	"time"

	"portwatch/internal/portschedule"
)

func TestNew_DefaultInterval(t *testing.T) {
	s := portschedule.New(30 * time.Second)
	if got := s.Interval(80); got != 30*time.Second {
		t.Fatalf("expected 30s default, got %v", got)
	}
}

func TestNew_ZeroDefaultFallsBack(t *testing.T) {
	s := portschedule.New(0)
	if got := s.Interval(443); got != 60*time.Second {
		t.Fatalf("expected 60s fallback, got %v", got)
	}
}

func TestAdd_ValidEntry(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	if err := s.Add(portschedule.Entry{MinPort: 1, MaxPort: 1024, Interval: 10 * time.Second}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", s.Len())
	}
}

func TestAdd_InvalidRangeReturnsError(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	err := s.Add(portschedule.Entry{MinPort: 1024, MaxPort: 80, Interval: 5 * time.Second})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestAdd_ZeroIntervalReturnsError(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	err := s.Add(portschedule.Entry{MinPort: 1, MaxPort: 100, Interval: 0})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestInterval_MatchesRange(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	_ = s.Add(portschedule.Entry{MinPort: 1, MaxPort: 1024, Interval: 15 * time.Second})
	if got := s.Interval(22); got != 15*time.Second {
		t.Fatalf("expected 15s for port 22, got %v", got)
	}
}

func TestInterval_OutsideRangeUsesDefault(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	_ = s.Add(portschedule.Entry{MinPort: 1, MaxPort: 1024, Interval: 15 * time.Second})
	if got := s.Interval(8080); got != 60*time.Second {
		t.Fatalf("expected 60s default for port 8080, got %v", got)
	}
}

func TestInterval_ShortestWins(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	_ = s.Add(portschedule.Entry{MinPort: 1, MaxPort: 1024, Interval: 20 * time.Second})
	_ = s.Add(portschedule.Entry{MinPort: 20, MaxPort: 25, Interval: 5 * time.Second})
	if got := s.Interval(22); got != 5*time.Second {
		t.Fatalf("expected 5s (shortest) for port 22, got %v", got)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	s := portschedule.New(60 * time.Second)
	_ = s.Add(portschedule.Entry{MinPort: 1, MaxPort: 100, Interval: 5 * time.Second})
	s.Reset()
	if s.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", s.Len())
	}
	if got := s.Interval(50); got != 60*time.Second {
		t.Fatalf("expected default after reset, got %v", got)
	}
}
