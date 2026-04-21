package portpolicy

import (
	"testing"
	"time"

	"portwatch/internal/snapshot"
)

func entry(port int, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func makeSnap(entries ...snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries, time.Now())
}

func TestNew_ReturnsEmptyPolicy(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("expected non-nil policy")
	}
}

func TestEvaluate_DefaultAllow(t *testing.T) {
	p := New()
	action, v := p.Evaluate(entry(8080, "tcp"))
	if action != Allow || v != nil {
		t.Fatalf("expected allow/nil, got %v/%v", action, v)
	}
}

func TestEvaluate_DenyMatchesPortAndProtocol(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "no-telnet", Port: 23, Protocol: "tcp", Action: Deny})
	action, v := p.Evaluate(entry(23, "tcp"))
	if action != Deny {
		t.Fatalf("expected Deny, got %v", action)
	}
	if v == nil || v.Rule != "no-telnet" {
		t.Fatalf("unexpected violation: %v", v)
	}
}

func TestEvaluate_ProtocolMismatch(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "no-telnet", Port: 23, Protocol: "tcp", Action: Deny})
	action, v := p.Evaluate(entry(23, "udp"))
	if action != Allow || v != nil {
		t.Fatalf("expected allow for protocol mismatch, got %v/%v", action, v)
	}
}

func TestEvaluate_EmptyProtocolMatchesAny(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "block-all-23", Port: 23, Protocol: "", Action: Deny})
	for _, proto := range []string{"tcp", "udp"} {
		action, v := p.Evaluate(entry(23, proto))
		if action != Deny || v == nil {
			t.Fatalf("proto=%s: expected Deny, got %v/%v", proto, action, v)
		}
	}
}

func TestEvaluate_FirstMatchWins(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "allow-22", Port: 22, Protocol: "tcp", Action: Allow})
	p.Add(Rule{Name: "deny-22", Port: 22, Protocol: "tcp", Action: Deny})
	action, v := p.Evaluate(entry(22, "tcp"))
	if action != Allow || v != nil {
		t.Fatalf("first-match-wins failed: got %v/%v", action, v)
	}
}

func TestCheck_ReturnsAllViolations(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "no-telnet", Port: 23, Protocol: "tcp", Action: Deny})
	p.Add(Rule{Name: "no-ftp", Port: 21, Protocol: "tcp", Action: Deny})

	snap := makeSnap(entry(22, "tcp"), entry(23, "tcp"), entry(21, "tcp"))
	violations := p.Check(snap)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestCheck_NilSnapshot(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "no-telnet", Port: 23, Protocol: "tcp", Action: Deny})
	if v := p.Check(nil); v != nil {
		t.Fatalf("expected nil for nil snapshot, got %v", v)
	}
}

func TestAction_String(t *testing.T) {
	if Allow.String() != "allow" {
		t.Fatalf("unexpected: %s", Allow.String())
	}
	if Deny.String() != "deny" {
		t.Fatalf("unexpected: %s", Deny.String())
	}
}
