package portmute

import (
	"testing"
	"time"
)

func TestNew_ReturnsEmptyMuter(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil Muter")
	}
	if len(m.Active()) != 0 {
		t.Fatalf("expected 0 active rules, got %d", len(m.Active()))
	}
}

func TestMute_IsMuted_BasicFlow(t *testing.T) {
	m := New()
	m.Mute(8080, "tcp", time.Hour, "test")
	if !m.IsMuted(8080, "tcp") {
		t.Fatal("expected port to be muted")
	}
}

func TestIsMuted_FalseForUnknownPort(t *testing.T) {
	m := New()
	if m.IsMuted(9999, "tcp") {
		t.Fatal("expected unknown port to not be muted")
	}
}

func TestMute_ProtocolDistinct(t *testing.T) {
	m := New()
	m.Mute(53, "tcp", time.Hour, "dns tcp mute")
	if !m.IsMuted(53, "tcp") {
		t.Fatal("expected tcp/53 to be muted")
	}
	if m.IsMuted(53, "udp") {
		t.Fatal("expected udp/53 to not be muted")
	}
}

func TestUnmute_RemovesRule(t *testing.T) {
	m := New()
	m.Mute(443, "tcp", time.Hour, "test")
	m.Unmute(443, "tcp")
	if m.IsMuted(443, "tcp") {
		t.Fatal("expected port to no longer be muted after Unmute")
	}
}

func TestUnmute_NoopForUnknownPort(t *testing.T) {
	m := New()
	m.Unmute(1234, "tcp") // should not panic
}

func TestIsMuted_ExpiresAfterDuration(t *testing.T) {
	m := New()
	now := time.Now()
	m.now = func() time.Time { return now }
	m.Mute(22, "tcp", 5*time.Minute, "ssh mute")

	// advance past expiry
	m.now = func() time.Time { return now.Add(6 * time.Minute) }
	if m.IsMuted(22, "tcp") {
		t.Fatal("expected mute to have expired")
	}
}

func TestActive_ReturnsCurrentRules(t *testing.T) {
	m := New()
	m.Mute(80, "tcp", time.Hour, "web")
	m.Mute(443, "tcp", time.Hour, "https")
	active := m.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active rules, got %d", len(active))
	}
}

func TestActive_PrunesExpiredRules(t *testing.T) {
	m := New()
	now := time.Now()
	m.now = func() time.Time { return now }
	m.Mute(80, "tcp", time.Minute, "short")
	m.Mute(443, "tcp", time.Hour, "long")

	m.now = func() time.Time { return now.Add(2 * time.Minute) }
	active := m.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active rule after expiry, got %d", len(active))
	}
	if active[0].Port != 443 {
		t.Fatalf("expected remaining rule for port 443, got %d", active[0].Port)
	}
}

func TestMute_ReasonStored(t *testing.T) {
	m := New()
	m.Mute(8080, "tcp", time.Hour, "maintenance")
	active := m.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active rule, got %d", len(active))
	}
	if active[0].Reason != "maintenance" {
		t.Fatalf("expected reason 'maintenance', got %q", active[0].Reason)
	}
}
