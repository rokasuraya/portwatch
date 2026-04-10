package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestNew_AllowsFirstEvent(t *testing.T) {
	th := throttle.New(time.Second)
	if !th.Allow("tcp:8080") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_SuppressesDuplicateWithinCooldown(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("tcp:8080") // prime
	if th.Allow("tcp:8080") {
		t.Fatal("expected duplicate within cooldown to be suppressed")
	}
}

func TestAllow_PermitsDifferentKeys(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("tcp:8080")
	if !th.Allow("udp:53") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestAllow_PermitsAfterCooldown(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow("tcp:9090")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("tcp:9090") {
		t.Fatal("expected event to be allowed after cooldown elapsed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("tcp:443")
	th.Reset("tcp:443")
	if !th.Allow("tcp:443") {
		t.Fatal("expected event to be allowed after reset")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow("tcp:22")
	time.Sleep(20 * time.Millisecond)
	th.Purge()
	// After purge the key is gone; Allow should return true again.
	if !th.Allow("tcp:22") {
		t.Fatal("expected purged key to be allowed again")
	}
}

func TestPurge_KeepsActiveEntries(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("tcp:80")
	th.Purge() // cooldown has not elapsed
	if th.Allow("tcp:80") {
		t.Fatal("expected active entry to survive purge")
	}
}
