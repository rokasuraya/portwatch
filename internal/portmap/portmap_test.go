package portmap

import (
	"testing"
)

func TestNew_Empty(t *testing.T) {
	pm := New()
	if pm.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", pm.Len())
	}
}

func TestSet_AddsEntry(t *testing.T) {
	pm := New()
	pm.Set(80, "tcp", "http", true)
	if pm.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", pm.Len())
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	pm := New()
	pm.Set(443, "tcp", "https", true)
	e, ok := pm.Get(443, "tcp")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Port != 443 || e.Protocol != "tcp" || e.Label != "https" || !e.Open {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestGet_MissingReturnsFalse(t *testing.T) {
	pm := New()
	_, ok := pm.Get(9999, "tcp")
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	pm := New()
	pm.Set(22, "tcp", "ssh", true)
	e, _ := pm.Get(22, "tcp")
	e.Label = "mutated"
	e2, _ := pm.Get(22, "tcp")
	if e2.Label == "mutated" {
		t.Fatal("Get should return a copy, not a reference")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	pm := New()
	pm.Set(8080, "tcp", "http-alt", true)
	pm.Delete(8080, "tcp")
	if pm.Len() != 0 {
		t.Fatalf("expected 0 entries after delete, got %d", pm.Len())
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	pm := New()
	pm.Set(80, "tcp", "http", true)
	pm.Set(53, "udp", "dns", true)
	all := pm.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestSet_OverwritesExistingEntry(t *testing.T) {
	pm := New()
	pm.Set(80, "tcp", "http", true)
	pm.Set(80, "tcp", "http-updated", false)
	e, ok := pm.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Label != "http-updated" || e.Open {
		t.Fatalf("expected overwritten entry, got %+v", e)
	}
}

func TestProtocol_DistinctKeys(t *testing.T) {
	pm := New()
	pm.Set(53, "tcp", "dns-tcp", true)
	pm.Set(53, "udp", "dns-udp", true)
	if pm.Len() != 2 {
		t.Fatalf("expected 2 entries for same port different proto, got %d", pm.Len())
	}
}
