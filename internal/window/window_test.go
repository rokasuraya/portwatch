package window_test

import (
	"testing"
	"time"

	"portwatch/internal/window"
)

func TestNew_ReturnsWindow(t *testing.T) {
	w := window.New(time.Second, 10)
	if w == nil {
		t.Fatal("expected non-nil Window")
	}
}

func TestCount_ZeroOnEmpty(t *testing.T) {
	w := window.New(time.Second, 10)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncreasesCount(t *testing.T) {
	w := window.New(time.Second, 10)
	w.Add(3)
	w.Add(2)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestCount_ExpiresOldBuckets(t *testing.T) {
	w := window.New(50*time.Millisecond, 10)
	w.Add(7)
	time.Sleep(80 * time.Millisecond)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCount_RetainsRecentBuckets(t *testing.T) {
	w := window.New(200*time.Millisecond, 10)
	w.Add(4)
	time.Sleep(20 * time.Millisecond)
	w.Add(6)
	if got := w.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	w := window.New(time.Second, 10)
	w.Add(5)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestNew_SizeZeroDefaultsToOne(t *testing.T) {
	// Should not panic with size=0.
	w := window.New(time.Second, 0)
	w.Add(1)
	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
