package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func entries(ports ...int) []scanner.Entry {
	out := make([]scanner.Entry, len(ports))
	for i, p := range ports {
		out[i] = scanner.Entry{Port: p, Protocol: "tcp"}
	}
	return out
}

func TestCompute_DeterministicForSameInput(t *testing.T) {
	a := fingerprint.Compute(entries(80, 443, 22))
	b := fingerprint.Compute(entries(80, 443, 22))
	if a != b {
		t.Fatalf("expected identical fingerprints, got %q vs %q", a, b)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := fingerprint.Compute(entries(22, 80, 443))
	b := fingerprint.Compute(entries(443, 22, 80))
	if a != b {
		t.Fatalf("expected order-independent fingerprints, got %q vs %q", a, b)
	}
}

func TestCompute_DiffersForDifferentPorts(t *testing.T) {
	a := fingerprint.Compute(entries(80))
	b := fingerprint.Compute(entries(8080))
	if a == b {
		t.Fatal("expected different fingerprints for different ports")
	}
}

func TestCompute_EmptySlice(t *testing.T) {
	fp := fingerprint.Compute(nil)
	if fp == "" {
		t.Fatal("expected non-empty fingerprint for empty slice")
	}
}

func TestNew_ZeroLast(t *testing.T) {
	tr := fingerprint.New()
	if tr.Last() != "" {
		t.Fatalf("expected empty last fingerprint, got %q", tr.Last())
	}
}

func TestChanged_TrueOnFirstCall(t *testing.T) {
	tr := fingerprint.New()
	fp := fingerprint.Compute(entries(80))
	if !tr.Changed(fp) {
		t.Fatal("expected Changed to return true on first call")
	}
}

func TestChanged_FalseWhenSame(t *testing.T) {
	tr := fingerprint.New()
	fp := fingerprint.Compute(entries(80, 443))
	tr.Changed(fp)
	if tr.Changed(fp) {
		t.Fatal("expected Changed to return false for identical fingerprint")
	}
}

func TestChanged_TrueAfterDiff(t *testing.T) {
	tr := fingerprint.New()
	tr.Changed(fingerprint.Compute(entries(80)))
	if !tr.Changed(fingerprint.Compute(entries(80, 443))) {
		t.Fatal("expected Changed to return true after port set changed")
	}
}

func TestReset_ClearsLast(t *testing.T) {
	tr := fingerprint.New()
	fp := fingerprint.Compute(entries(22))
	tr.Changed(fp)
	tr.Reset()
	if tr.Last() != "" {
		t.Fatalf("expected empty last after Reset, got %q", tr.Last())
	}
	if !tr.Changed(fp) {
		t.Fatal("expected Changed to return true after Reset")
	}
}
