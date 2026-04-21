package portrank_test

import (
	"testing"

	"portwatch/internal/portrank"
	"portwatch/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_ReturnsRanker(t *testing.T) {
	r := portrank.New()
	if r == nil {
		t.Fatal("expected non-nil Ranker")
	}
}

func TestRank_NilSnapshot(t *testing.T) {
	r := portrank.New(portrank.PrivilegedPortScorer())
	if got := r.Rank(nil); got != nil {
		t.Fatalf("expected nil result, got %v", got)
	}
}

func TestRank_EmptySnapshot(t *testing.T) {
	r := portrank.New(portrank.PrivilegedPortScorer())
	snap := makeSnap(nil)
	got := r.Rank(snap)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(got))
	}
}

func TestRank_SortedDescending(t *testing.T) {
	r := portrank.New(
		portrank.PrivilegedPortScorer(),
		portrank.WellKnownRiskyPortScorer(),
	)
	snap := makeSnap([]snapshot.Entry{
		{Port: 8080, Protocol: "tcp"},
		{Port: 22, Protocol: "tcp"},
		{Port: 80, Protocol: "tcp"},
	})
	got := r.Rank(snap)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	// Port 22 should rank first (privileged + risky)
	if got[0].Port != 22 {
		t.Errorf("expected port 22 first, got %d", got[0].Port)
	}
	for i := 1; i < len(got); i++ {
		if got[i].Score > got[i-1].Score {
			t.Errorf("entries not sorted: index %d score %.2f > index %d score %.2f",
				i, got[i].Score, i-1, got[i-1].Score)
		}
	}
}

func TestRank_ReasonsPopulated(t *testing.T) {
	r := portrank.New(portrank.WellKnownRiskyPortScorer())
	snap := makeSnap([]snapshot.Entry{
		{Port: 3389, Protocol: "tcp"},
	})
	got := r.Rank(snap)
	if len(got) == 0 {
		t.Fatal("expected one entry")
	}
	if len(got[0].Reasons) == 0 {
		t.Error("expected at least one reason for port 3389")
	}
}

func TestAddScorer_AppendsDynamically(t *testing.T) {
	r := portrank.New()
	snap := makeSnap([]snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
	})

	before := r.Rank(snap)
	if before[0].Score != 0 {
		t.Fatalf("expected zero score before scorer added, got %.2f", before[0].Score)
	}

	r.AddScorer(portrank.WellKnownRiskyPortScorer())
	after := r.Rank(snap)
	if after[0].Score == 0 {
		t.Error("expected non-zero score after scorer added")
	}
}

func TestProtocolScorer_UDPPenalty(t *testing.T) {
	r := portrank.New(portrank.ProtocolScorer())
	snap := makeSnap([]snapshot.Entry{
		{Port: 53, Protocol: "udp"},
		{Port: 53, Protocol: "tcp"},
	})
	got := r.Rank(snap)
	if got[0].Protocol != "udp" {
		t.Errorf("expected udp to rank higher, got %s first", got[0].Protocol)
	}
}
