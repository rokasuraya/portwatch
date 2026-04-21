package portcap_test

import (
	"testing"

	"portwatch/internal/portcap"
	"portwatch/internal/snapshot"
)

func makeEntry(port uint16, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_ReturnsCap(t *testing.T) {
	c := portcap.New(10, nil)
	if c == nil {
		t.Fatal("expected non-nil PortCap")
	}
}

func TestCheck_UnderLimit(t *testing.T) {
	c := portcap.New(5, nil)
	snap := makeSnap([]snapshot.Entry{
		makeEntry(80, "tcp"),
		makeEntry(443, "tcp"),
	})
	violations := c.Check(snap)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_ExceedsLimit(t *testing.T) {
	c := portcap.New(2, nil)
	snap := makeSnap([]snapshot.Entry{
		makeEntry(80, "tcp"),
		makeEntry(443, "tcp"),
		makeEntry(8080, "tcp"),
	})
	violations := c.Check(snap)
	if len(violations) == 0 {
		t.Fatal("expected violations when over cap")
	}
}

func TestCheck_ExactLimit(t *testing.T) {
	c := portcap.New(3, nil)
	snap := makeSnap([]snapshot.Entry{
		makeEntry(22, "tcp"),
		makeEntry(80, "tcp"),
		makeEntry(443, "tcp"),
	})
	violations := c.Check(snap)
	if len(violations) != 0 {
		t.Fatalf("expected no violations at exact limit, got %d", len(violations))
	}
}

func TestCheck_NilSnapshot(t *testing.T) {
	c := portcap.New(5, nil)
	violations := c.Check(nil)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for nil snapshot, got %d", len(violations))
	}
}

func TestCheck_ZeroMaxAlwaysViolates(t *testing.T) {
	c := portcap.New(0, nil)
	snap := makeSnap([]snapshot.Entry{
		makeEntry(80, "tcp"),
	})
	violations := c.Check(snap)
	if len(violations) == 0 {
		t.Fatal("expected violation when max is 0 and ports exist")
	}
}
