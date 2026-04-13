package presencecheck_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/presencecheck"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_DefaultsToStdout(t *testing.T) {
	checker := presencecheck.New(nil, nil)
	if checker == nil {
		t.Fatal("expected non-nil checker")
	}
}

func TestCheck_AllPresent(t *testing.T) {
	required := []snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
		{Port: 80, Protocol: "tcp"},
	}
	snap := makeSnap(required)
	checker := presencecheck.New(required, nil)
	results := checker.Check(snap)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Present {
			t.Errorf("expected port %d/%s to be present", r.Port, r.Protocol)
		}
	}
}

func TestCheck_MissingPort(t *testing.T) {
	required := []snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	snap := makeSnap([]snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
	})
	checker := presencecheck.New(required, nil)
	results := checker.Check(snap)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	var missing int
	for _, r := range results {
		if !r.Present {
			missing++
		}
	}
	if missing != 1 {
		t.Errorf("expected 1 missing, got %d", missing)
	}
}

func TestCheck_ProtocolDistinct(t *testing.T) {
	required := []snapshot.Entry{
		{Port: 53, Protocol: "tcp"},
		{Port: 53, Protocol: "udp"},
	}
	snap := makeSnap([]snapshot.Entry{
		{Port: 53, Protocol: "tcp"},
	})
	checker := presencecheck.New(required, nil)
	results := checker.Check(snap)

	presentCount := 0
	for _, r := range results {
		if r.Present {
			presentCount++
		}
	}
	if presentCount != 1 {
		t.Errorf("expected only tcp/53 present, got %d present", presentCount)
	}
}

func TestReport_WritesMissingPorts(t *testing.T) {
	var buf bytes.Buffer
	required := []snapshot.Entry{
		{Port: 8080, Protocol: "tcp"},
	}
	snap := makeSnap(nil)
	checker := presencecheck.New(required, &buf)
	results := checker.Check(snap)
	checker.Report(results)

	if !strings.Contains(buf.String(), "MISSING") {
		t.Errorf("expected MISSING in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "8080") {
		t.Errorf("expected port 8080 in output, got: %s", buf.String())
	}
}

func TestReport_SilentWhenAllPresent(t *testing.T) {
	var buf bytes.Buffer
	required := []snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
	}
	snap := makeSnap(required)
	checker := presencecheck.New(required, &buf)
	results := checker.Check(snap)
	checker.Report(results)

	if buf.Len() != 0 {
		t.Errorf("expected no output when all ports present, got: %s", buf.String())
	}
}
