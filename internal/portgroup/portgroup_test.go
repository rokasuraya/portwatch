package portgroup

import (
	"sort"
	"testing"
)

func TestNew_ReturnsEmptyRegistry(t *testing.T) {
	r := New()
	if len(r.Names()) != 0 {
		t.Fatalf("expected no groups, got %d", len(r.Names()))
	}
}

func TestDefine_AddsGroup(t *testing.T) {
	r := New()
	if err := r.Define("web", []string{"80/tcp", "443/tcp"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := r.Names()
	if len(names) != 1 || names[0] != "web" {
		t.Fatalf("expected [web], got %v", names)
	}
}

func TestDefine_RejectsEmptyEntry(t *testing.T) {
	r := New()
	if err := r.Define("bad", []string{"22/tcp", ""}); err == nil {
		t.Fatal("expected error for empty entry")
	}
}

func TestContains_ReturnsTrueForMember(t *testing.T) {
	r := New()
	_ = r.Define("db", []string{"5432/tcp", "3306/tcp"})
	if !r.Contains("db", "5432/tcp") {
		t.Error("expected 5432/tcp to be in group db")
	}
}

func TestContains_ReturnsFalseForNonMember(t *testing.T) {
	r := New()
	_ = r.Define("db", []string{"5432/tcp"})
	if r.Contains("db", "22/tcp") {
		t.Error("expected 22/tcp not to be in group db")
	}
}

func TestContains_ReturnsFalseForUnknownGroup(t *testing.T) {
	r := New()
	if r.Contains("missing", "80/tcp") {
		t.Error("expected false for unknown group")
	}
}

func TestMembers_ReturnsAllEntries(t *testing.T) {
	r := New()
	_ = r.Define("ssh", []string{"22/tcp"})
	m := r.Members("ssh")
	if len(m) != 1 || m[0] != "22/tcp" {
		t.Fatalf("unexpected members: %v", m)
	}
}

func TestMembers_NilForUnknownGroup(t *testing.T) {
	r := New()
	if r.Members("nope") != nil {
		t.Error("expected nil for unknown group")
	}
}

func TestRemove_DeletesGroup(t *testing.T) {
	r := New()
	_ = r.Define("tmp", []string{"9999/tcp"})
	r.Remove("tmp")
	if r.Contains("tmp", "9999/tcp") {
		t.Error("group should have been removed")
	}
}

func TestDefine_ReplacesExistingGroup(t *testing.T) {
	r := New()
	_ = r.Define("web", []string{"80/tcp"})
	_ = r.Define("web", []string{"443/tcp"})
	m := r.Members("web")
	sort.Strings(m)
	if len(m) != 1 || m[0] != "443/tcp" {
		t.Fatalf("expected only 443/tcp after replace, got %v", m)
	}
}
