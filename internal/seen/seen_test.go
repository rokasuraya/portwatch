package seen_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/seen"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(entries ...snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_EmptyLedger(t *testing.T) {
	l := seen.New()
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", l.Len())
	}
}

func TestObserve_NilSnapshot(t *testing.T) {
	l := seen.New()
	l.Observe(nil) // must not panic
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries after nil observe")
	}
}

func TestObserve_RecordsFirstSeen(t *testing.T) {
	l := seen.New()
	before := time.Now()
	snap := makeSnap(snapshot.Entry{Protocol: "tcp", Port: 80})
	l.Observe(snap)

	e, ok := l.Lookup(snapshot.Entry{Protocol: "tcp", Port: 80})
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.FirstSeen.Before(before) {
		t.Errorf("FirstSeen %v is before test start %v", e.FirstSeen, before)
	}
	if e.Count != 1 {
		t.Errorf("expected Count=1, got %d", e.Count)
	}
}

func TestObserve_IncrementsCount(t *testing.T) {
	l := seen.New()
	snap := makeSnap(snapshot.Entry{Protocol: "tcp", Port: 443})
	l.Observe(snap)
	l.Observe(snap)
	l.Observe(snap)

	e, ok := l.Lookup(snapshot.Entry{Protocol: "tcp", Port: 443})
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Count != 3 {
		t.Errorf("expected Count=3, got %d", e.Count)
	}
}

func TestObserve_ProtocolDistinct(t *testing.T) {
	l := seen.New()
	l.Observe(makeSnap(snapshot.Entry{Protocol: "tcp", Port: 53}))
	l.Observe(makeSnap(snapshot.Entry{Protocol: "udp", Port: 53}))

	if l.Len() != 2 {
		t.Errorf("expected 2 distinct entries, got %d", l.Len())
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	l := seen.New()
	_, ok := l.Lookup(snapshot.Entry{Protocol: "tcp", Port: 9999})
	if ok {
		t.Fatal("expected no entry for unknown port")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	l := seen.New()
	l.Observe(makeSnap(
		snapshot.Entry{Protocol: "tcp", Port: 22},
		snapshot.Entry{Protocol: "tcp", Port: 80},
	))
	if l.Len() == 0 {
		t.Fatal("expected entries before reset")
	}
	l.Reset()
	if l.Len() != 0 {
		t.Errorf("expected 0 entries after reset, got %d", l.Len())
	}
}
