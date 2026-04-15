package portname

import "testing"

func TestNew_NilOverrides(t *testing.T) {
	m := New(nil)
	if m == nil {
		t.Fatal("expected non-nil Mapper")
	}
}

func TestLookup_WellKnownPort(t *testing.T) {
	m := New(nil)
	tests := []struct {
		proto string
		port  int
		want  string
	}{
		{"tcp", 22, "ssh"},
		{"tcp", 80, "http"},
		{"tcp", 443, "https"},
		{"udp", 53, "dns"},
		{"tcp", 3306, "mysql"},
	}
	for _, tt := range tests {
		got := m.Lookup(tt.proto, tt.port)
		if got != tt.want {
			t.Errorf("Lookup(%s, %d) = %q; want %q", tt.proto, tt.port, got, tt.want)
		}
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	m := New(nil)
	got := m.Lookup("tcp", 9999)
	if got != "" {
		t.Errorf("expected empty string for unknown port, got %q", got)
	}
}

func TestLookup_CustomOverridesBuiltIn(t *testing.T) {
	m := New(map[string]string{"tcp:80": "my-http"})
	got := m.Lookup("tcp", 80)
	if got != "my-http" {
		t.Errorf("expected custom override 'my-http', got %q", got)
	}
}

func TestRegister_AddsMapping(t *testing.T) {
	m := New(nil)
	m.Register("tcp", 9000, "myservice")
	got := m.Lookup("tcp", 9000)
	if got != "myservice" {
		t.Errorf("expected 'myservice', got %q", got)
	}
}

func TestRegister_OverridesExisting(t *testing.T) {
	m := New(nil)
	m.Register("tcp", 22, "custom-ssh")
	got := m.Lookup("tcp", 22)
	if got != "custom-ssh" {
		t.Errorf("expected 'custom-ssh', got %q", got)
	}
}

func TestLookupWithFallback_Known(t *testing.T) {
	m := New(nil)
	got := m.LookupWithFallback("tcp", 443)
	if got != "https" {
		t.Errorf("expected 'https', got %q", got)
	}
}

func TestLookupWithFallback_Unknown(t *testing.T) {
	m := New(nil)
	got := m.LookupWithFallback("tcp", 12345)
	if got != "port/12345" {
		t.Errorf("expected 'port/12345', got %q", got)
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	m := New(nil)
	tcp := m.Lookup("tcp", 53)
	udp := m.Lookup("udp", 53)
	if tcp != "dns" || udp != "dns" {
		t.Errorf("both tcp:53 and udp:53 should resolve to 'dns', got tcp=%q udp=%q", tcp, udp)
	}
	// ensure tcp-only port is not found via udp
	got := m.Lookup("udp", 22)
	if got != "" {
		t.Errorf("udp:22 should be empty, got %q", got)
	}
}
