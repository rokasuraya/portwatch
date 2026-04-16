package portquota_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portquota"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(ports []int) *snapshot.Snapshot {
	entries := make([]scanner.Entry, len(ports))
	for i, p := range ports {
		entries[i] = scanner.Entry{Port: p, Proto: "tcp"}
	}
	return snapshot.New(entries, time.Now())
}

func TestNew_DefaultsToStderr(t *testing.T) {
	q := portquota.New(10, nil)
	if q == nil {
		t.Fatal("expected non-nil Quota")
	}
	if q.Max() != 10 {
		t.Fatalf("expected max 10, got %d", q.Max())
	}
}

func TestCheck_UnderLimit(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(5, &buf)
	snap := makeSnap([]int{80, 443, 8080})
	if q.Check(snap) {
		t.Error("expected false for count under limit")
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestCheck_ExceedsLimit(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(2, &buf)
	snap := makeSnap([]int{80, 443, 8080})
	if !q.Check(snap) {
		t.Error("expected true for count exceeding limit")
	}
	if !strings.Contains(buf.String(), "exceeds limit") {
		t.Errorf("expected warning in output, got: %s", buf.String())
	}
}

func TestCheck_WarnOnce(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(1, &buf)
	snap := makeSnap([]int{80, 443})
	q.Check(snap)
	q.Check(snap)
	count := strings.Count(buf.String(), "exceeds limit")
	if count != 1 {
		t.Errorf("expected warning once, got %d times", count)
	}
}

func TestCheck_RecoveryResetsSuppression(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(2, &buf)
	q.Check(makeSnap([]int{80, 443, 8080})) // breach
	q.Check(makeSnap([]int{80}))            // recover
	buf.Reset()
	q.Check(makeSnap([]int{80, 443, 8080})) // breach again
	if !strings.Contains(buf.String(), "exceeds limit") {
		t.Error("expected warning after recovery")
	}
}

func TestCheck_NilSnapshot(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(5, &buf)
	if q.Check(nil) {
		t.Error("expected false for nil snapshot")
	}
}

func TestSetMax_UpdatesLimit(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(10, &buf)
	q.SetMax(2)
	if q.Max() != 2 {
		t.Fatalf("expected max 2, got %d", q.Max())
	}
	snap := makeSnap([]int{80, 443, 8080})
	if !q.Check(snap) {
		t.Error("expected breach after SetMax")
	}
}
