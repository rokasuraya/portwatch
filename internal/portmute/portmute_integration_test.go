package portmute_test

import (
	"testing"
	"time"

	"portwatch/internal/portmute"
)

func TestIntegration_MuteWindowSuppressesAndRecovers(t *testing.T) {
	m := portmute.New()

	// Nothing muted initially.
	if m.IsMuted(8080, "tcp") {
		t.Fatal("expected no mute before any rule added")
	}

	// Apply a very short mute.
	m.Mute(8080, "tcp", 50*time.Millisecond, "integration test")
	if !m.IsMuted(8080, "tcp") {
		t.Fatal("expected port to be muted immediately after rule added")
	}

	// Wait for expiry.
	time.Sleep(100 * time.Millisecond)
	if m.IsMuted(8080, "tcp") {
		t.Fatal("expected mute to have expired")
	}
}

func TestIntegration_MultiplePortsMutedIndependently(t *testing.T) {
	m := portmute.New()
	m.Mute(22, "tcp", time.Hour, "ssh maintenance")
	m.Mute(3306, "tcp", time.Hour, "db maintenance")
	m.Mute(5432, "tcp", time.Hour, "pg maintenance")

	for _, port := range []int{22, 3306, 5432} {
		if !m.IsMuted(port, "tcp") {
			t.Errorf("expected port %d to be muted", port)
		}
	}

	m.Unmute(3306, "tcp")
	if m.IsMuted(3306, "tcp") {
		t.Error("expected port 3306 to be unmuted")
	}
	if !m.IsMuted(22, "tcp") {
		t.Error("expected port 22 to remain muted")
	}
	if !m.IsMuted(5432, "tcp") {
		t.Error("expected port 5432 to remain muted")
	}

	active := m.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active rules, got %d", len(active))
	}
}
