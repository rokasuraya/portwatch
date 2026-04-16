package portlock

import (
	"bytes"
	"testing"
)

func TestNew_ReturnsLocker(t *testing.T) {
	l := New(nil)
	if l == nil {
		t.Fatal("expected non-nil Locker")
	}
	if l.Len() != 0 {
		t.Fatalf("expected 0 locked ports, got %d", l.Len())
	}
}

func TestLock_AddsEntry(t *testing.T) {
	l := New(new(bytes.Buffer))
	l.Lock(22, "tcp", "ssh baseline")
	if !l.IsLocked(22, "tcp") {
		t.Fatal("expected port 22/tcp to be locked")
	}
}

func TestIsLocked_FalseForUnknownPort(t *testing.T) {
	l := New(new(bytes.Buffer))
	if l.IsLocked(80, "tcp") {
		t.Fatal("expected port 80/tcp to not be locked")
	}
}

func TestUnlock_RemovesEntry(t *testing.T) {
	l := New(new(bytes.Buffer))
	l.Lock(443, "tcp", "")
	l.Unlock(443, "tcp")
	if l.IsLocked(443, "tcp") {
		t.Fatal("expected port 443/tcp to be unlocked after Unlock")
	}
}

func TestUnlock_NoopForUnknownPort(t *testing.T) {
	l := New(new(bytes.Buffer))
	l.Unlock(9999, "udp") // should not panic
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", l.Len())
	}
}

func TestLock_ProtocolDistinct(t *testing.T) {
	l := New(new(bytes.Buffer))
	l.Lock(53, "tcp", "dns tcp")
	l.Lock(53, "udp", "dns udp")
	if l.Len() != 2 {
		t.Fatalf("expected 2 locked entries, got %d", l.Len())
	}
	if !l.IsLocked(53, "tcp") || !l.IsLocked(53, "udp") {
		t.Fatal("expected both 53/tcp and 53/udp to be locked")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	l := New(new(bytes.Buffer))
	l.Lock(8080, "tcp", "dev server")
	entries := l.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].Protocol != "tcp" {
		t.Fatalf("unexpected entry: %+v", entries[0])
	}
}

func TestLen_TracksChanges(t *testing.T) {
	l := New(new(bytes.Buffer))
	l.Lock(1, "tcp", "")
	l.Lock(2, "tcp", "")
	if l.Len() != 2 {
		t.Fatalf("expected 2, got %d", l.Len())
	}
	l.Unlock(1, "tcp")
	if l.Len() != 1 {
		t.Fatalf("expected 1, got %d", l.Len())
	}
}
