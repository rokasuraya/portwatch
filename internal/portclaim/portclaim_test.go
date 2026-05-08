package portclaim

import (
	"testing"
)

func TestNew_ReturnsEmptyRegistry(t *testing.T) {
	r := New()
	if r.Len() != 0 {
		t.Fatalf("expected 0 claims, got %d", r.Len())
	}
}

func TestRegister_AddsEntry(t *testing.T) {
	r := New()
	r.Register(Claim{Port: 80, Protocol: "tcp", Owner: "nginx", PID: 1234})
	if r.Len() != 1 {
		t.Fatalf("expected 1 claim, got %d", r.Len())
	}
}

func TestLookup_ReturnsRegisteredClaim(t *testing.T) {
	r := New()
	want := Claim{Port: 443, Protocol: "tcp", Owner: "caddy", PID: 5678}
	r.Register(want)

	got, ok := r.Lookup(443, "tcp")
	if !ok {
		t.Fatal("expected claim to be found")
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestLookup_MissingReturnsFalse(t *testing.T) {
	r := New()
	_, ok := r.Lookup(9999, "tcp")
	if ok {
		t.Fatal("expected no claim for unknown port")
	}
}

func TestRegister_OverwritesExistingClaim(t *testing.T) {
	r := New()
	r.Register(Claim{Port: 22, Protocol: "tcp", Owner: "sshd", PID: 100})
	r.Register(Claim{Port: 22, Protocol: "tcp", Owner: "dropbear", PID: 200})

	got, ok := r.Lookup(22, "tcp")
	if !ok {
		t.Fatal("expected claim to be found")
	}
	if got.Owner != "dropbear" {
		t.Fatalf("expected owner dropbear, got %s", got.Owner)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 claim after overwrite, got %d", r.Len())
	}
}

func TestProtocolDistinct(t *testing.T) {
	r := New()
	r.Register(Claim{Port: 53, Protocol: "tcp", Owner: "bind", PID: 300})
	r.Register(Claim{Port: 53, Protocol: "udp", Owner: "bind", PID: 300})

	if r.Len() != 2 {
		t.Fatalf("expected 2 claims for tcp/udp, got %d", r.Len())
	}
}

func TestRelease_RemovesClaim(t *testing.T) {
	r := New()
	r.Register(Claim{Port: 8080, Protocol: "tcp", Owner: "app", PID: 999})
	r.Release(8080, "tcp")

	_, ok := r.Lookup(8080, "tcp")
	if ok {
		t.Fatal("expected claim to be removed after Release")
	}
	if r.Len() != 0 {
		t.Fatalf("expected 0 claims after release, got %d", r.Len())
	}
}

func TestRelease_NoopForUnknownPort(t *testing.T) {
	r := New()
	r.Release(1234, "tcp") // should not panic
	if r.Len() != 0 {
		t.Fatalf("expected 0 claims, got %d", r.Len())
	}
}
