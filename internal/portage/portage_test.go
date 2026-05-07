package portage

import (
	"testing"
	"time"

	"portwatch/internal/snapshot"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeSnap(ports []int, proto string) *snapshot.Snapshot {
	entries := make([]snapshot.Entry, len(ports))
	for i, p := range ports {
		entries[i] = snapshot.Entry{Port: p, Protocol: proto}
	}
	return snapshot.New(entries)
}

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_ReturnsTracker(t *testing.T) {
	tr := New(nil)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker")
	}
}

func TestObserve_RecordsFirstSeen(t *testing.T) {
	tr := New(fixedNow(epoch))
	tr.Observe(makeSnap([]int{80}, "tcp"))
	e, ok := tr.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if !e.FirstSeen.Equal(epoch) {
		t.Errorf("FirstSeen = %v, want %v", e.FirstSeen, epoch)
	}
	if e.SeenCount != 1 {
		t.Errorf("SeenCount = %d, want 1", e.SeenCount)
	}
}

func TestObserve_DoesNotResetFirstSeen(t *testing.T) {
	tr := New(fixedNow(epoch))
	tr.Observe(makeSnap([]int{443}, "tcp"))

	later := epoch.Add(5 * time.Minute)
	tr.now = fixedNow(later)
	tr.Observe(makeSnap([]int{443}, "tcp"))

	e, ok := tr.Get(443, "tcp")
	if !ok {
		t.Fatal("expected entry")
	}
	if !e.FirstSeen.Equal(epoch) {
		t.Errorf("FirstSeen should not change: got %v", e.FirstSeen)
	}
	if !e.LastSeen.Equal(later) {
		t.Errorf("LastSeen = %v, want %v", e.LastSeen, later)
	}
	if e.SeenCount != 2 {
		t.Errorf("SeenCount = %d, want 2", e.SeenCount)
	}
}

func TestObserve_RemovesClosedPort(t *testing.T) {
	tr := New(fixedNow(epoch))
	tr.Observe(makeSnap([]int{22, 80}, "tcp"))
	tr.Observe(makeSnap([]int{80}, "tcp"))

	if _, ok := tr.Get(22, "tcp"); ok {
		t.Error("port 22 should have been removed")
	}
	if _, ok := tr.Get(80, "tcp"); !ok {
		t.Error("port 80 should still be present")
	}
}

func TestObserve_NilSnapshotNoOp(t *testing.T) {
	tr := New(fixedNow(epoch))
	tr.Observe(makeSnap([]int{8080}, "tcp"))
	tr.Observe(nil)
	if _, ok := tr.Get(8080, "tcp"); !ok {
		t.Error("nil snapshot should not clear entries")
	}
}

func TestEntry_Age(t *testing.T) {
	e := Entry{FirstSeen: epoch}
	got := e.Age(epoch.Add(10 * time.Minute))
	if got != 10*time.Minute {
		t.Errorf("Age = %v, want 10m", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := New(fixedNow(epoch))
	tr.Observe(makeSnap([]int{22, 80, 443}, "tcp"))
	all := tr.All()
	if len(all) != 3 {
		t.Errorf("len(All) = %d, want 3", len(all))
	}
}
