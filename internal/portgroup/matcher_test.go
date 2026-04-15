package portgroup

import (
	"sort"
	"testing"
)

func setup() *Registry {
	r := New()
	_ = r.Define("web", []string{"80/tcp", "443/tcp"})
	_ = r.Define("ssh", []string{"22/tcp"})
	_ = r.Define("db", []string{"5432/tcp", "3306/tcp"})
	return r
}

func TestNewMatcher_ReturnsMatcher(t *testing.T) {
	m := NewMatcher(setup())
	if m == nil {
		t.Fatal("expected non-nil matcher")
	}
}

func TestMatch_ReturnsMatchingGroups(t *testing.T) {
	m := NewMatcher(setup())
	res := m.Match(443, "tcp")
	if res.Key != "443/tcp" {
		t.Errorf("unexpected key: %s", res.Key)
	}
	if len(res.Groups) != 1 || res.Groups[0] != "web" {
		t.Errorf("expected [web], got %v", res.Groups)
	}
}

func TestMatch_ReturnsEmptyWhenNoMatch(t *testing.T) {
	m := NewMatcher(setup())
	res := m.Match(9999, "tcp")
	if len(res.Groups) != 0 {
		t.Errorf("expected no groups, got %v", res.Groups)
	}
}

func TestMatch_MultipleGroups(t *testing.T) {
	r := New()
	_ = r.Define("all", []string{"22/tcp", "80/tcp"})
	_ = r.Define("ssh", []string{"22/tcp"})
	m := NewMatcher(r)
	res := m.Match(22, "tcp")
	sort.Strings(res.Groups)
	if len(res.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %v", res.Groups)
	}
	if res.Groups[0] != "all" || res.Groups[1] != "ssh" {
		t.Errorf("unexpected groups: %v", res.Groups)
	}
}

func TestAnyMatch_TrueForKnownPort(t *testing.T) {
	m := NewMatcher(setup())
	if !m.AnyMatch(22, "tcp") {
		t.Error("expected true for 22/tcp")
	}
}

func TestAnyMatch_FalseForUnknownPort(t *testing.T) {
	m := NewMatcher(setup())
	if m.AnyMatch(12345, "tcp") {
		t.Error("expected false for unknown port")
	}
}

func TestInGroup_TrueForMember(t *testing.T) {
	m := NewMatcher(setup())
	if !m.InGroup("db", 5432, "tcp") {
		t.Error("expected 5432/tcp to be in db")
	}
}

func TestInGroup_FalseForNonMember(t *testing.T) {
	m := NewMatcher(setup())
	if m.InGroup("ssh", 443, "tcp") {
		t.Error("expected 443/tcp not to be in ssh")
	}
}
