package backoff

import (
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	b := New(0, 0, 0)
	if b.base != 100*time.Millisecond {
		t.Fatalf("expected default base 100ms, got %v", b.base)
	}
	if b.max != 30*time.Second {
		t.Fatalf("expected default max 30s, got %v", b.max)
	}
	if b.factor != 2.0 {
		t.Fatalf("expected default factor 2.0, got %v", b.factor)
	}
}

func TestNext_GrowsExponentially(t *testing.T) {
	b := New(100*time.Millisecond, 10*time.Second, 2.0)

	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}

	for i, want := range expected {
		got := b.Next()
		if got != want {
			t.Fatalf("attempt %d: expected %v, got %v", i, want, got)
		}
	}
}

func TestNext_CapsAtMax(t *testing.T) {
	b := New(100*time.Millisecond, 300*time.Millisecond, 2.0)

	for i := 0; i < 10; i++ {
		d := b.Next()
		if d > 300*time.Millisecond {
			t.Fatalf("attempt %d: delay %v exceeded max 300ms", i, d)
		}
	}
}

func TestAttempts_TracksCount(t *testing.T) {
	b := New(10*time.Millisecond, time.Second, 2.0)
	for i := 0; i < 5; i++ {
		b.Next()
	}
	if b.Attempts() != 5 {
		t.Fatalf("expected 5 attempts, got %d", b.Attempts())
	}
}

func TestReset_ClearsAttempts(t *testing.T) {
	b := New(10*time.Millisecond, time.Second, 2.0)
	b.Next()
	b.Next()
	b.Reset()

	if b.Attempts() != 0 {
		t.Fatalf("expected 0 attempts after reset, got %d", b.Attempts())
	}
	// First call after reset should return base delay.
	if got := b.Next(); got != 10*time.Millisecond {
		t.Fatalf("expected base delay after reset, got %v", got)
	}
}

func TestNext_ConcurrentSafe(t *testing.T) {
	b := New(time.Millisecond, time.Second, 2.0)
	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			b.Next()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
	if b.Attempts() != 20 {
		t.Fatalf("expected 20 concurrent attempts, got %d", b.Attempts())
	}
}
