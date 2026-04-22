package portpulse

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port uint16, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func makeDiff(opened, closed []snapshot.Entry) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	tr := New(time.Minute, nil)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
	if tr.out == nil {
		t.Fatal("expected default writer")
	}
}

func TestObserve_CountsEvents(t *testing.T) {
	tr := New(time.Minute, &bytes.Buffer{})
	tr.Observe(makeDiff(
		[]snapshot.Entry{makeEntry(80, "tcp"), makeEntry(443, "tcp")},
		[]snapshot.Entry{makeEntry(8080, "tcp")},
	))
	if got := tr.Rate(); got != 3 {
		t.Fatalf("expected rate 3, got %d", got)
	}
}

func TestObserve_EmptyDiffNoChange(t *testing.T) {
	tr := New(time.Minute, &bytes.Buffer{})
	tr.Observe(makeDiff(nil, nil))
	if got := tr.Rate(); got != 0 {
		t.Fatalf("expected rate 0, got %d", got)
	}
}

func TestRate_PrunesExpiredEvents(t *testing.T) {
	tr := New(50*time.Millisecond, &bytes.Buffer{})
	tr.Observe(makeDiff(
		[]snapshot.Entry{makeEntry(22, "tcp")},
		nil,
	))
	if tr.Rate() != 1 {
		t.Fatal("expected 1 before expiry")
	}
	time.Sleep(80 * time.Millisecond)
	if got := tr.Rate(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestRate_AccumulatesMultipleObservations(t *testing.T) {
	tr := New(time.Minute, &bytes.Buffer{})
	tr.Observe(makeDiff([]snapshot.Entry{makeEntry(80, "tcp")}, nil))
	tr.Observe(makeDiff(nil, []snapshot.Entry{makeEntry(443, "tcp")}))
	tr.Observe(makeDiff([]snapshot.Entry{makeEntry(22, "tcp"), makeEntry(3306, "tcp")}, nil))
	if got := tr.Rate(); got != 4 {
		t.Fatalf("expected rate 4, got %d", got)
	}
}

func TestReport_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	tr := New(time.Minute, &buf)
	tr.Observe(makeDiff([]snapshot.Entry{makeEntry(80, "tcp")}, nil))
	tr.Report()
	if !strings.Contains(buf.String(), "portpulse:") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "1 change events") {
		t.Fatalf("expected count in output: %q", buf.String())
	}
}

func TestReport_ZeroRateWhenEmpty(t *testing.T) {
	var buf bytes.Buffer
	tr := New(time.Minute, &buf)
	tr.Report()
	if !strings.Contains(buf.String(), "0 change events") {
		t.Fatalf("expected zero count: %q", buf.String())
	}
}
