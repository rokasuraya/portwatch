package jitter_test

import (
	"testing"
	"time"

	"portwatch/internal/jitter"
)

func TestNew_ReturnsSource(t *testing.T) {
	s := jitter.New()
	if s == nil {
		t.Fatal("expected non-nil Source")
	}
}

func TestApply_ZeroMaxReturnsBase(t *testing.T) {
	s := jitter.NewWithSeed(42)
	base := 5 * time.Second
	got := s.Apply(base, 0)
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_NegativeMaxReturnsBase(t *testing.T) {
	s := jitter.NewWithSeed(42)
	base := 5 * time.Second
	got := s.Apply(base, -1*time.Second)
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_AddsOffsetWithinRange(t *testing.T) {
	s := jitter.NewWithSeed(1)
	base := 1 * time.Second
	max := 500 * time.Millisecond
	for i := 0; i < 100; i++ {
		result := s.Apply(base, max)
		if result < base {
			t.Fatalf("result %v less than base %v", result, base)
		}
		if result >= base+max {
			t.Fatalf("result %v exceeds base+max %v", result, base+max)
		}
	}
}

func TestSpread_ZeroMaxReturnsZero(t *testing.T) {
	s := jitter.NewWithSeed(42)
	if got := s.Spread(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestSpread_ReturnsValueWithinRange(t *testing.T) {
	s := jitter.NewWithSeed(7)
	max := 200 * time.Millisecond
	for i := 0; i < 100; i++ {
		v := s.Spread(max)
		if v < 0 || v >= max {
			t.Fatalf("spread %v out of range [0, %v)", v, max)
		}
	}
}

func TestClamp_DoesNotExceedCeiling(t *testing.T) {
	s := jitter.NewWithSeed(99)
	base := 900 * time.Millisecond
	max := 500 * time.Millisecond
	ceiling := 1 * time.Second
	for i := 0; i < 200; i++ {
		result := s.Clamp(base, max, ceiling)
		if result > ceiling {
			t.Fatalf("result %v exceeds ceiling %v", result, ceiling)
		}
		if result < base {
			t.Fatalf("result %v less than base %v", result, base)
		}
	}
}

func TestApply_DeterministicWithFixedSeed(t *testing.T) {
	s1 := jitter.NewWithSeed(123)
	s2 := jitter.NewWithSeed(123)
	base := 2 * time.Second
	max := 1 * time.Second
	for i := 0; i < 20; i++ {
		if s1.Apply(base, max) != s2.Apply(base, max) {
			t.Fatal("expected identical sequences from same seed")
		}
	}
}
