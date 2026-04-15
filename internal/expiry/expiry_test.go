package expiry

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(ports []int, proto string) *snapshot.Snapshot {
	entries := make([]snapshot.Entry, len(ports))
	for i, p := range ports {
		entries[i] = snapshot.Entry{Port: p, Protocol: proto}
	}
	return snapshot.New(entries)
}

func TestNew_DefaultsToStdout(t *testing.T) {
	c := New(time.Minute, nil)
	if c == nil {
		t.Fatal("expected non-nil Checker")
	}
	if c.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestObserve_RecordsNewPorts(t *testing.T) {
	var buf bytes.Buffer
	c := New(time.Hour, &buf)
	snap := makeSnap([]int{80, 443}, "tcp")
	c.Observe(snap)

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(c.records))
	}
}

func TestObserve_RemovesClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	c := New(time.Hour, &buf)
	c.Observe(makeSnap([]int{80, 443}, "tcp"))
	c.Observe(makeSnap([]int{80}, "tcp")) // 443 closed

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) != 1 {
		t.Fatalf("expected 1 record after close, got %d", len(c.records))
	}
	if _, ok := c.records[portKey(80, "tcp")]; !ok {
		t.Error("expected port 80 to remain")
	}
}

func TestCheck_NoWarningWhenFresh(t *testing.T) {
	var buf bytes.Buffer
	c := New(time.Hour, &buf)
	c.Observe(makeSnap([]int{22}, "tcp"))

	count := c.Check()
	if count != 0 {
		t.Fatalf("expected 0 expired, got %d", count)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestCheck_WarnsWhenExpired(t *testing.T) {
	var buf bytes.Buffer
	c := New(time.Millisecond, &buf)
	c.Observe(makeSnap([]int{8080}, "tcp"))

	// backdate the record so it appears old
	c.mu.Lock()
	k := portKey(8080, "tcp")
	rec := c.records[k]
	rec.FirstSeen = time.Now().Add(-time.Hour)
	c.records[k] = rec
	c.mu.Unlock()

	count := c.Check()
	if count != 1 {
		t.Fatalf("expected 1 expired entry, got %d", count)
	}
	if !strings.Contains(buf.String(), "8080") {
		t.Errorf("expected port 8080 in output, got: %s", buf.String())
	}
}

func TestCheck_ProtocolDistinct(t *testing.T) {
	var buf bytes.Buffer
	c := New(time.Millisecond, &buf)

	tcpEntries := []snapshot.Entry{{Port: 53, Protocol: "tcp"}}
	udpEntries := []snapshot.Entry{{Port: 53, Protocol: "udp"}}
	allEntries := append(tcpEntries, udpEntries...)
	snap := snapshot.New(allEntries)
	c.Observe(snap)

	c.mu.Lock()
	for k := range c.records {
		rec := c.records[k]
		rec.FirstSeen = time.Now().Add(-time.Hour)
		c.records[k] = rec
	}
	c.mu.Unlock()

	count := c.Check()
	if count != 2 {
		t.Fatalf("expected 2 expired (tcp+udp), got %d", count)
	}
}
