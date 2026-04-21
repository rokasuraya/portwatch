package portpolicy

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"portwatch/internal/snapshot"
)

func goodSnap() (*snapshot.Snapshot, error) {
	return makeSnap(entry(23, "tcp")), nil
}

func badSnap() (*snapshot.Snapshot, error) {
	return nil, errors.New("scan failed")
}

func TestNewEnforcer_DefaultsToStderr(t *testing.T) {
	p := New()
	e := NewEnforcer(p, goodSnap, time.Second, nil)
	if e == nil {
		t.Fatal("expected non-nil enforcer")
	}
}

func TestEnforcer_WritesViolation(t *testing.T) {
	p := New()
	p.Add(Rule{Name: "no-telnet", Port: 23, Protocol: "tcp", Action: Deny})

	var buf bytes.Buffer
	e := NewEnforcer(p, goodSnap, time.Second, &buf)
	e.enforce()

	if !strings.Contains(buf.String(), "no-telnet") {
		t.Fatalf("expected violation output, got: %q", buf.String())
	}
}

func TestEnforcer_SnapshotError_WritesMessage(t *testing.T) {
	p := New()
	var buf bytes.Buffer
	e := NewEnforcer(p, badSnap, time.Second, &buf)
	e.enforce()

	if !strings.Contains(buf.String(), "snapshot error") {
		t.Fatalf("expected error message, got: %q", buf.String())
	}
}

func TestEnforcer_Run_CancelsCleanly(t *testing.T) {
	p := New()
	var buf bytes.Buffer
	e := NewEnforcer(p, goodSnap, 10*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()

	err := e.Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
