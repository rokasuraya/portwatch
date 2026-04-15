package portmap

import (
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port int, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func TestNewUpdater_NilLabels(t *testing.T) {
	pm := New()
	u := NewUpdater(pm, nil)
	u.Apply([]snapshot.Entry{makeEntry(80, "tcp")}, nil)
	e, ok := pm.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Label != "" {
		t.Fatalf("expected empty label, got %q", e.Label)
	}
}

func TestApply_OpenedAddsEntries(t *testing.T) {
	pm := New()
	u := NewUpdater(pm, func(port int, proto string) string {
		if port == 22 {
			return "ssh"
		}
		return ""
	})
	u.Apply([]snapshot.Entry{makeEntry(22, "tcp"), makeEntry(80, "tcp")}, nil)
	if pm.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", pm.Len())
	}
	e, _ := pm.Get(22, "tcp")
	if e.Label != "ssh" {
		t.Fatalf("expected label ssh, got %q", e.Label)
	}
}

func TestApply_ClosedRemovesEntries(t *testing.T) {
	pm := New()
	pm.Set(8080, "tcp", "http-alt", true)
	u := NewUpdater(pm, nil)
	u.Apply(nil, []snapshot.Entry{makeEntry(8080, "tcp")})
	if pm.Len() != 0 {
		t.Fatalf("expected 0 entries after close, got %d", pm.Len())
	}
}

func TestApply_MixedDiff(t *testing.T) {
	pm := New()
	pm.Set(443, "tcp", "https", true)
	u := NewUpdater(pm, nil)
	u.Apply(
		[]snapshot.Entry{makeEntry(22, "tcp")},
		[]snapshot.Entry{makeEntry(443, "tcp")},
	)
	if pm.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", pm.Len())
	}
	if _, ok := pm.Get(22, "tcp"); !ok {
		t.Fatal("expected port 22 to be present")
	}
	if _, ok := pm.Get(443, "tcp"); ok {
		t.Fatal("expected port 443 to be removed")
	}
}
