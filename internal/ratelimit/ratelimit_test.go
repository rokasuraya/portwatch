package ratelimit_test

import (
	"testing"
	"time"

	"portwatch/internal/ratelimit"
)

func TestNew_DefaultsBurstToOne(t *testing.T) {
	l := ratelimit.New(time.Second, 0)
	if l.Remaining() != 1 {
		t.Fatalf("expected burst 1, got %d", l.Remaining())
	}
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(time.Second, 1)
	if !l.Allow() {
		t.Fatal("expected first Allow() to return true")
	}
}

func TestAllow_SecondCallBlockedWithinInterval(t *testing.T) {
	l := ratelimit.New(time.Second, 1)
	l.Allow()
	if l.Allow() {
		t.Fatal("expected second Allow() within interval to return false")
	}
}

func TestAllow_BurstPermitsMultipleCalls(t *testing.T) {
	l := ratelimit.New(time.Second, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() to return true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow() to return false after burst exhausted")
	}
}

func TestAllow_RefillsAfterInterval(t *testing.T) {
	l := ratelimit.New(50*time.Millisecond, 1)
	l.Allow() // consume token
	time.Sleep(60 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected Allow() to return true after interval elapsed")
	}
}

func TestReset_RestoresTokens(t *testing.T) {
	l := ratelimit.New(time.Second, 2)
	l.Allow()
	l.Allow()
	if l.Remaining() != 0 {
		t.Fatalf("expected 0 tokens after exhaustion, got %d", l.Remaining())
	}
	l.Reset()
	if l.Remaining() != 2 {
		t.Fatalf("expected 2 tokens after reset, got %d", l.Remaining())
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	l := ratelimit.New(time.Second, 3)
	if l.Remaining() != 3 {
		t.Fatalf("expected 3 remaining, got %d", l.Remaining())
	}
	l.Allow()
	if l.Remaining() != 2 {
		t.Fatalf("expected 2 remaining after one Allow, got %d", l.Remaining())
	}
}
