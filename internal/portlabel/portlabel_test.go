package portlabel

import (
	"testing"
)

func TestNew_NilOverrides(t *testing.T) {
	l := New(nil)
	if l == nil {
		t.Fatal("expected non-nil Labeler")
	}
	if l.Len() == 0 {
		t.Fatal("expected built-in labels to be loaded")
	}
}

func TestNew_ExtraOverridesBuiltIn(t *testing.T) {
	overrides := map[string]Label{
		"80/tcp": {Name: "CustomHTTP", Category: "custom"},
	}
	l := New(overrides)
	lbl, ok := l.Lookup(80, "tcp")
	if !ok {
		t.Fatal("expected ok=true for overridden port")
	}
	if lbl.Name != "CustomHTTP" {
		t.Fatalf("expected CustomHTTP, got %s", lbl.Name)
	}
}

func TestLookup_WellKnownPort(t *testing.T) {
	l := New(nil)
	lbl, ok := l.Lookup(22, "tcp")
	if !ok {
		t.Fatal("expected ok=true for SSH")
	}
	if lbl.Name != "SSH" {
		t.Fatalf("expected SSH, got %s", lbl.Name)
	}
	if lbl.Category != "remote-access" {
		t.Fatalf("expected remote-access, got %s", lbl.Category)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	l := New(nil)
	lbl, ok := l.Lookup(9999, "tcp")
	if ok {
		t.Fatal("expected ok=false for unknown port")
	}
	if lbl.Category != "unknown" {
		t.Fatalf("expected unknown category, got %s", lbl.Category)
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	l := New(nil)
	_, tcpOK := l.Lookup(53, "tcp")
	_, udpOK := l.Lookup(53, "udp")
	if !tcpOK || !udpOK {
		t.Fatal("expected both tcp and udp DNS entries")
	}
	_, sctp := l.Lookup(53, "sctp")
	if sctp {
		t.Fatal("expected ok=false for unknown protocol variant")
	}
}

func TestRegister_AddsNewEntry(t *testing.T) {
	l := New(nil)
	before := l.Len()
	l.Register(12345, "tcp", Label{Name: "MyApp", Category: "app"})
	if l.Len() != before+1 {
		t.Fatalf("expected Len to grow by 1")
	}
	lbl, ok := l.Lookup(12345, "tcp")
	if !ok {
		t.Fatal("expected registered entry to be found")
	}
	if lbl.Name != "MyApp" {
		t.Fatalf("expected MyApp, got %s", lbl.Name)
	}
}

func TestRegister_ReplacesExisting(t *testing.T) {
	l := New(nil)
	l.Register(22, "tcp", Label{Name: "SecureShell", Category: "custom"})
	lbl, _ := l.Lookup(22, "tcp")
	if lbl.Name != "SecureShell" {
		t.Fatalf("expected SecureShell, got %s", lbl.Name)
	}
}
