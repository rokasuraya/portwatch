package suppress_test

import (
	"testing"
	"time"

	"portwatch/internal/suppress"
)

func TestNew_ReturnsEmptySuppressor(t *testing.T) {
	s := suppress.New()
	if s == nil {
		t.Fatal("expected non-nil Suppressor")
	}
	if got := len(s.List()); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestIsSuppressed_FalseForUnknownPort(t *testing.T) {
	s := suppress.New()
	if s.IsSuppressed(8080, "tcp") {
		t.Fatal("expected false for unknown port")
	}
}

func TestAdd_SuppressesPort(t *testing.T) {
	s := suppress.New()
	s.Add(22, "tcp", "known SSH", time.Time{})
	if !s.IsSuppressed(22, "tcp") {
		t.Fatal("expected port 22/tcp to be suppressed")
	}
}

func TestAdd_ProtocolDistinct(t *testing.T) {
	s := suppress.New()
	s.Add(53, "tcp", "DNS tcp", time.Time{})
	if s.IsSuppressed(53, "udp") {
		t.Fatal("udp should not be suppressed when only tcp was added")
	}
}

func TestRemove_ClearsRule(t *testing.T) {
	s := suppress.New()
	s.Add(443, "tcp", "HTTPS", time.Time{})
	s.Remove(443, "tcp")
	if s.IsSuppressed(443, "tcp") {
		t.Fatal("expected port 443/tcp to no longer be suppressed")
	}
}

func TestIsSuppressed_ExpiredRuleReturnsFalse(t *testing.T) {
	s := suppress.New()
	past := time.Now().Add(-1 * time.Second)
	s.Add(9000, "tcp", "temp", past)
	if s.IsSuppressed(9000, "tcp") {
		t.Fatal("expected expired rule to be treated as absent")
	}
}

func TestList_ExcludesExpiredEntries(t *testing.T) {
	s := suppress.New()
	s.Add(80, "tcp", "HTTP", time.Time{})
	s.Add(9999, "udp", "expired", time.Now().Add(-time.Second))
	entries := s.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 active entry, got %d", len(entries))
	}
	if entries[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", entries[0].Port)
	}
}

func TestList_ReturnsAllActiveEntries(t *testing.T) {
	s := suppress.New()
	future := time.Now().Add(time.Hour)
	s.Add(22, "tcp", "SSH", time.Time{})
	s.Add(3306, "tcp", "MySQL", future)
	if got := len(s.List()); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}
