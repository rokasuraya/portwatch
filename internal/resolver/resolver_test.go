package resolver_test

import (
	"testing"

	"github.com/user/portwatch/internal/resolver"
)

func TestNew_NilExtra(t *testing.T) {
	r := resolver.New(nil)
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestLookup_WellKnownPort(t *testing.T) {
	r := resolver.New(nil)
	got := r.Lookup(22, "tcp")
	if got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	r := resolver.New(nil)
	got := r.Lookup(9999, "tcp")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestLookup_ExtraOverridesBuiltIn(t *testing.T) {
	r := resolver.New(map[string]string{"tcp/22": "custom-ssh"})
	got := r.Lookup(22, "tcp")
	if got != "custom-ssh" {
		t.Fatalf("expected custom-ssh, got %q", got)
	}
}

func TestLookup_ExtraOnlyPort(t *testing.T) {
	r := resolver.New(map[string]string{"tcp/12345": "myapp"})
	got := r.Lookup(12345, "tcp")
	if got != "myapp" {
		t.Fatalf("expected myapp, got %q", got)
	}
}

func TestRegister_AddsMapping(t *testing.T) {
	r := resolver.New(nil)
	r.Register(9200, "tcp", "elasticsearch")
	got := r.Lookup(9200, "tcp")
	if got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestRegister_OverridesExisting(t *testing.T) {
	r := resolver.New(nil)
	r.Register(80, "tcp", "my-http")
	got := r.Lookup(80, "tcp")
	if got != "my-http" {
		t.Fatalf("expected my-http, got %q", got)
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	r := resolver.New(nil)
	tcp := r.Lookup(53, "tcp")
	udp := r.Lookup(53, "udp")
	if tcp != "dns" || udp != "dns" {
		t.Fatalf("expected both dns, got tcp=%q udp=%q", tcp, udp)
	}
	got := r.Lookup(53, "sctp")
	if got != "" {
		t.Fatalf("expected empty for sctp/53, got %q", got)
	}
}
