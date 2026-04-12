package digest_test

import (
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/snapshot"
)

func entries(pairs ...any) []snapshot.Entry {
	var out []snapshot.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, snapshot.Entry{
			Proto: pairs[i].(string),
			Port:  pairs[i+1].(int),
		})
	}
	return out
}

func TestCompute_DeterministicForSameInput(t *testing.T) {
	e := entries("tcp", 80, "tcp", 443)
	if digest.Compute(e) != digest.Compute(e) {
		t.Fatal("expected same digest for identical input")
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := digest.Compute(entries("tcp", 80, "tcp", 443))
	b := digest.Compute(entries("tcp", 443, "tcp", 80))
	if a != b {
		t.Fatalf("expected order-independent digest, got %s vs %s", a, b)
	}
}

func TestCompute_DiffersForDifferentPorts(t *testing.T) {
	a := digest.Compute(entries("tcp", 80))
	b := digest.Compute(entries("tcp", 8080))
	if digest.Equal(a, b) {
		t.Fatal("expected different digests for different ports")
	}
}

func TestCompute_EmptySlice(t *testing.T) {
	d := digest.Compute(nil)
	if d.String() == "" {
		t.Fatal("expected non-empty digest for empty input")
	}
}

func TestTracker_ChangedOnFirstCall(t *testing.T) {
	tr := digest.New()
	_, changed := tr.Changed(entries("tcp", 22))
	if !changed {
		t.Fatal("expected changed=true on first call")
	}
}

func TestTracker_UnchangedOnRepeat(t *testing.T) {
	tr := digest.New()
	e := entries("tcp", 22)
	tr.Changed(e)
	_, changed := tr.Changed(e)
	if changed {
		t.Fatal("expected changed=false when snapshot is identical")
	}
}

func TestTracker_ChangedAfterPortAdded(t *testing.T) {
	tr := digest.New()
	tr.Changed(entries("tcp", 22))
	_, changed := tr.Changed(entries("tcp", 22, "tcp", 80))
	if !changed {
		t.Fatal("expected changed=true after port added")
	}
}

func TestTracker_ResetCausesTrueOnNext(t *testing.T) {
	tr := digest.New()
	e := entries("tcp", 22)
	tr.Changed(e)
	tr.Reset()
	_, changed := tr.Changed(e)
	if !changed {
		t.Fatal("expected changed=true after Reset")
	}
}

func TestTracker_LastReturnsCurrentDigest(t *testing.T) {
	tr := digest.New()
	e := entries("udp", 53)
	d, _ := tr.Changed(e)
	if tr.Last() != d {
		t.Fatalf("Last() = %s, want %s", tr.Last(), d)
	}
}
