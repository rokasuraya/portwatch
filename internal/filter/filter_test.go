package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestNew_EmptyRules(t *testing.T) {
	f := filter.New(nil)
	if f == nil {
		t.Fatal("expected non-nil Filter")
	}
}

func TestAllow_PermitsUnknownPort(t *testing.T) {
	f := filter.New(nil)
	if !f.Allow(8080, "tcp") {
		t.Error("expected port 8080/tcp to be allowed when no rules defined")
	}
}

func TestAllow_SuppressesMatchingRule(t *testing.T) {
	rules := []filter.Rule{
		{Port: 22, Protocol: "tcp", Comment: "SSH"},
	}
	f := filter.New(rules)
	if f.Allow(22, "tcp") {
		t.Error("expected port 22/tcp to be suppressed")
	}
}

func TestAllow_ProtocolMismatch(t *testing.T) {
	rules := []filter.Rule{
		{Port: 53, Protocol: "tcp"},
	}
	f := filter.New(rules)
	if !f.Allow(53, "udp") {
		t.Error("expected port 53/udp to be allowed; rule only covers tcp")
	}
}

func TestApply_FiltersEntries(t *testing.T) {
	rules := []filter.Rule{
		{Port: 22, Protocol: "tcp"},
		{Port: 80, Protocol: "tcp"},
	}
	f := filter.New(rules)

	input := []string{"22/tcp", "80/tcp", "443/tcp", "8080/tcp"}
	got := f.Apply(input)

	want := []string{"443/tcp", "8080/tcp"}
	if len(got) != len(want) {
		t.Fatalf("Apply() returned %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("Apply()[%d] = %q, want %q", i, got[i], v)
		}
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f := filter.New(nil)
	got := f.Apply([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestApply_KeepsUnparseable(t *testing.T) {
	f := filter.New(nil)
	input := []string{"not-a-port"}
	got := f.Apply(input)
	if len(got) != 1 || got[0] != "not-a-port" {
		t.Errorf("expected unparseable entry to be kept, got %v", got)
	}
}
