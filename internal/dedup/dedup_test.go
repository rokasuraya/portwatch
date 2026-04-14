package dedup_test

import (
	"testing"

	"github.com/user/portwatch/internal/dedup"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(ports ...int) *snapshot.Snapshot {
	entries := make([]scanner.Entry, len(ports))
	for i, p := range ports {
		entries[i] = scanner.Entry{Port: p, Protocol: "tcp"}
	}
	return snapshot.New(entries)
}

func TestNew_ReturnsDeduplicator(t *testing.T) {
	d := dedup.New()
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestAccept_FirstSnapshotAlwaysAccepted(t *testing.T) {
	d := dedup.New()
	snap := makeSnap(80, 443)
	if !d.Accept(snap) {
		t.Error("first snapshot should always be accepted")
	}
}

func TestAccept_IdenticalSnapshotRejected(t *testing.T) {
	d := dedup.New()
	snap := makeSnap(80, 443)
	d.Accept(snap)

	if d.Accept(makeSnap(80, 443)) {
		t.Error("identical snapshot should be rejected")
	}
}

func TestAccept_DifferentSnapshotAccepted(t *testing.T) {
	d := dedup.New()
	d.Accept(makeSnap(80))

	if !d.Accept(makeSnap(80, 443)) {
		t.Error("different snapshot should be accepted")
	}
}

func TestAccept_NilReturnsFalse(t *testing.T) {
	d := dedup.New()
	if d.Accept(nil) {
		t.Error("nil snapshot should return false")
	}
}

func TestReset_AllowsReacceptance(t *testing.T) {
	d := dedup.New()
	snap := makeSnap(80, 443)
	d.Accept(snap)
	d.Reset()

	if !d.Accept(makeSnap(80, 443)) {
		t.Error("after Reset, identical snapshot should be accepted again")
	}
}

func TestAccept_EmptySnapshotDeduplicated(t *testing.T) {
	d := dedup.New()
	d.Accept(makeSnap())

	if d.Accept(makeSnap()) {
		t.Error("two consecutive empty snapshots should deduplicate")
	}
}
