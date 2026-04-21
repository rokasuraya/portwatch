package portevict

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port int, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func TestNew_ReturnsEvictor(t *testing.T) {
	e := New(time.Minute)
	if e == nil {
		t.Fatal("expected non-nil Evictor")
	}
	if e.Len() != 0 {
		t.Fatalf("expected 0 evictions, got %d", e.Len())
	}
}

func TestEvict_AddsEntry(t *testing.T) {
	e := New(time.Minute)
	entry := makeEntry(8080, "tcp")
	e.Evict(entry)
	if e.Len() != 1 {
		t.Fatalf("expected 1, got %d", e.Len())
	}
}

func TestIsEvicted_TrueWithinQuiet(t *testing.T) {
	e := New(time.Minute)
	entry := makeEntry(443, "tcp")
	e.Evict(entry)
	if !e.IsEvicted(entry) {
		t.Fatal("expected entry to be evicted")
	}
}

func TestIsEvicted_FalseForUnknown(t *testing.T) {
	e := New(time.Minute)
	if e.IsEvicted(makeEntry(22, "tcp")) {
		t.Fatal("expected false for unknown port")
	}
}

func TestIsEvicted_FalseAfterQuietExpires(t *testing.T) {
	e := New(50 * time.Millisecond)
	entry := makeEntry(9000, "udp")
	e.Evict(entry)
	// advance internal clock past quiet period
	base := e.now()
	e.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	if e.IsEvicted(entry) {
		t.Fatal("expected false after quiet period")
	}
	if e.Len() != 0 {
		t.Fatalf("expected eviction to be pruned, got %d", e.Len())
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	e := New(time.Minute)
	entry := makeEntry(80, "tcp")
	e.Evict(entry)
	e.Clear(entry)
	if e.IsEvicted(entry) {
		t.Fatal("expected entry cleared")
	}
	if e.Len() != 0 {
		t.Fatalf("expected 0, got %d", e.Len())
	}
}

func TestProtocolDistinct(t *testing.T) {
	e := New(time.Minute)
	tcp := makeEntry(53, "tcp")
	udp := makeEntry(53, "udp")
	e.Evict(tcp)
	if e.IsEvicted(udp) {
		t.Fatal("udp should not be evicted when only tcp was evicted")
	}
}
