package portfreq_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portfreq"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries, time.Now())
}

func TestNew_ReturnsEmptyTracker(t *testing.T) {
	tr := portfreq.New()
	_, ok := tr.Get(80, "tcp")
	if ok {
		t.Fatal("expected empty tracker")
	}
}

func TestObserve_IncrementsCount(t *testing.T) {
	tr := portfreq.New()
	snap := makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}})
	tr.Observe(snap)
	e, ok := tr.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry after observe")
	}
	if e.Count != 1 {
		t.Fatalf("expected count 1, got %d", e.Count)
	}
}

func TestObserve_AccumulatesAcrossScans(t *testing.T) {
	tr := portfreq.New()
	snap := makeSnap([]snapshot.Entry{{Port: 443, Protocol: "tcp"}})
	tr.Observe(snap)
	tr.Observe(snap)
	tr.Observe(snap)
	e, _ := tr.Get(443, "tcp")
	if e.Count != 3 {
		t.Fatalf("expected count 3, got %d", e.Count)
	}
}

func TestObserve_NilSnapshotNoOp(t *testing.T) {
	tr := portfreq.New()
	tr.Observe(nil) // must not panic
}

func TestObserve_ProtocolDistinct(t *testing.T) {
	tr := portfreq.New()
	snap := makeSnap([]snapshot.Entry{
		{Port: 53, Protocol: "tcp"},
		{Port: 53, Protocol: "udp"},
	})
	tr.Observe(snap)
	tcp, _ := tr.Get(53, "tcp")
	udp, _ := tr.Get(53, "udp")
	if tcp.Count != 1 || udp.Count != 1 {
		t.Fatalf("expected both counts 1, got tcp=%d udp=%d", tcp.Count, udp.Count)
	}
}

func TestTop_ReturnsDescendingOrder(t *testing.T) {
	tr := portfreq.New()
	one := makeSnap([]snapshot.Entry{{Port: 22, Protocol: "tcp"}})
	three := makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}})
	for i := 0; i < 3; i++ {
		tr.Observe(three)
	}
	tr.Observe(one)

	top := tr.Top(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Port != 80 {
		t.Fatalf("expected port 80 first, got %d", top[0].Port)
	}
}

func TestTop_ZeroReturnsAll(t *testing.T) {
	tr := portfreq.New()
	snap := makeSnap([]snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	})
	tr.Observe(snap)
	if got := len(tr.Top(0)); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestReset_ClearsAllCounts(t *testing.T) {
	tr := portfreq.New()
	tr.Observe(makeSnap([]snapshot.Entry{{Port: 8080, Protocol: "tcp"}}))
	tr.Reset()
	_, ok := tr.Get(8080, "tcp")
	if ok {
		t.Fatal("expected empty tracker after reset")
	}
}
