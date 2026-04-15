package cooldown

import (
	"testing"
	"time"
)

func TestNew_ReturnsCooldown(t *testing.T) {
	c := New(time.Second)
	if c == nil {
		t.Fatal("expected non-nil Cooldown")
	}
	if c.Len() != 0 {
		t.Fatalf("expected 0 tracked keys, got %d", c.Len())
	}
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	c := New(time.Second)
	if !c.Allow("port:80:tcp") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallSuppressedWithinPeriod(t *testing.T) {
	c := New(time.Hour)
	c.Allow("port:80:tcp")
	if c.Allow("port:80:tcp") {
		t.Fatal("expected second call within period to be suppressed")
	}
}

func TestAllow_PermitsDifferentKeys(t *testing.T) {
	c := New(time.Hour)
	c.Allow("port:80:tcp")
	if !c.Allow("port:443:tcp") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestAllow_PermitsAfterPeriodExpires(t *testing.T) {
	c := New(10 * time.Millisecond)
	now := time.Now()
	c.now = func() time.Time { return now }
	c.Allow("port:22:tcp")

	// advance clock past the period
	c.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	if !c.Allow("port:22:tcp") {
		t.Fatal("expected call after period to be allowed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	c := New(time.Hour)
	c.Allow("port:80:tcp")
	c.Reset("port:80:tcp")
	if !c.Allow("port:80:tcp") {
		t.Fatal("expected Allow to pass after Reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	c := New(time.Hour)
	c.Allow("a")
	c.Allow("b")
	c.Allow("a") // duplicate — should not increase count
	if c.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", c.Len())
	}
}

func TestReset_UnknownKeyIsNoop(t *testing.T) {
	c := New(time.Hour)
	// should not panic
	c.Reset("nonexistent")
	if c.Len() != 0 {
		t.Fatalf("expected 0 keys after reset of unknown key, got %d", c.Len())
	}
}
