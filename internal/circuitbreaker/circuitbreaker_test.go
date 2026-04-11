package circuitbreaker

import (
	"testing"
	"time"
)

func TestNew_DefaultsApplied(t *testing.T) {
	cb := New(0, 0)
	if cb.threshold != 3 {
		t.Errorf("expected default threshold 3, got %d", cb.threshold)
	}
	if cb.resetTimeout != 30*time.Second {
		t.Errorf("expected default resetTimeout 30s, got %v", cb.resetTimeout)
	}
}

func TestAllow_ClosedStatePermits(t *testing.T) {
	cb := New(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	cb := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.CurrentState() != StateOpen {
		t.Fatalf("expected StateOpen, got %v", cb.CurrentState())
	}
}

func TestAllow_RejectsWhenOpen(t *testing.T) {
	cb := New(1, time.Minute)
	cb.RecordFailure()
	if err := cb.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_TransitionsToHalfOpenAfterTimeout(t *testing.T) {
	cb := New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if cb.CurrentState() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", cb.CurrentState())
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	cb := New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow() // transition to half-open
	cb.RecordSuccess()
	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", cb.CurrentState())
	}
}

func TestRecordSuccess_ResetsFailureCount(t *testing.T) {
	cb := New(3, time.Minute)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	if cb.failures != 0 {
		t.Fatalf("expected failures=0 after success, got %d", cb.failures)
	}
}

func TestAllow_DoesNotOpenBeforeThreshold(t *testing.T) {
	cb := New(3, time.Minute)
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected StateClosed before threshold, got %v", cb.CurrentState())
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil before threshold, got %v", err)
	}
}
