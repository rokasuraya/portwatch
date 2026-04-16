package portage

import (
	"testing"
	"time"

	"github.com/username/portwatch/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_ReturnsTracker(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
	if tr.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", tr.Len())
	}
}

func TestObserve_RecordsFirstSeen(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := New()
	tr.now = fixedNow(base)

	snap := makeSnap([]snapshot.Entry{{Port: 80, Proto: "tcp"}})
	tr.Observe(snap)

	if tr.Len() != 1 {
		t.Fatalf("expected 1 tracked port, got %d", tr.Len())
	}
	e, ok := tr.Age(80, "tcp")
	if !ok {
		t.Fatal("expected entry for port 80/tcp")
	}
	if !e.FirstSeen.Equal(base) {
		t.Fatalf("expected FirstSeen %v, got %v", base, e.FirstSeen)
	}
}

func TestObserve_DoesNotResetFirstSeen(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	later := base.Add(5 * time.Minute)
	tr := New()
	tr.now = fixedNow(base)

	snap := makeSnap([]snapshot.Entry{{Port: 443, Proto: "tcp"}})
	tr.Observe(snap)

	tr.now = fixedNow(later)
	tr.Observe(snap) // same port, should not reset

	e, ok := tr.Age(443, "tcp")
	if !ok {
		t.Fatal("expected entry for port 443/tcp")
	}
	if !e.FirstSeen.Equal(base) {
		t.Fatalf("FirstSeen should not change: got %v", e.FirstSeen)
	}
	if e.Age < 5*time.Minute {
		t.Fatalf("expected age >= 5m, got %v", e.Age)
	}
}

func TestObserve_RemovesClosedPorts(t *testing.T) {
	tr := New()
	tr.now = fixedNow(time.Now())

	snap := makeSnap([]snapshot.Entry{{Port: 22, Proto: "tcp"}})
	tr.Observe(snap)
	if tr.Len() != 1 {
		t.Fatalf("expected 1, got %d", tr.Len())
	}

	tr.Observe(makeSnap(nil)) // port closed
	if tr.Len() != 0 {
		t.Fatalf("expected 0 after close, got %d", tr.Len())
	}
	_, ok := tr.Age(22, "tcp")
	if ok {
		t.Fatal("expected no entry for closed port")
	}
}

func TestObserve_NilSnapshot(t *testing.T) {
	tr := New()
	tr.Observe(nil) // must not panic
	if tr.Len() != 0 {
		t.Fatalf("expected 0, got %d", tr.Len())
	}
}

func TestAge_UnknownPort(t *testing.T) {
	tr := New()
	_, ok := tr.Age(9999, "tcp")
	if ok {
		t.Fatal("expected false for unknown port")
	}
}
