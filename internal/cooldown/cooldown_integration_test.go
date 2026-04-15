package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

func TestIntegration_BurstOfEventsOnlyFirstPasses(t *testing.T) {
	cd := cooldown.New(50 * time.Millisecond)

	keys := []string{"port:80:tcp", "port:443:tcp", "port:22:tcp"}
	for _, k := range keys {
		if !cd.Allow(k) {
			t.Fatalf("first call for key %q should be allowed", k)
		}
	}

	// second pass — all should be suppressed
	for _, k := range keys {
		if cd.Allow(k) {
			t.Fatalf("second call for key %q within period should be suppressed", k)
		}
	}

	// wait for period to expire
	time.Sleep(60 * time.Millisecond)

	// third pass — all should be allowed again
	for _, k := range keys {
		if !cd.Allow(k) {
			t.Fatalf("call after period for key %q should be allowed", k)
		}
	}
}

func TestIntegration_ResetAllowsImmediateRetrigger(t *testing.T) {
	cd := cooldown.New(time.Hour)

	cd.Allow("port:8080:tcp")
	cd.Reset("port:8080:tcp")

	if !cd.Allow("port:8080:tcp") {
		t.Fatal("expected Allow to pass immediately after Reset")
	}

	if cd.Len() != 1 {
		t.Fatalf("expected 1 tracked key, got %d", cd.Len())
	}
}
