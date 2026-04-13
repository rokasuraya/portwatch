package presencecheck_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/presencecheck"
	"github.com/user/portwatch/internal/snapshot"
)

// TestIntegration_CheckAndReport_FullFlow exercises Check followed by Report
// to confirm the two methods compose correctly end-to-end.
func TestIntegration_CheckAndReport_FullFlow(t *testing.T) {
	var buf bytes.Buffer

	required := []snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	// Only 22 and 80 are up.
	snap := makeSnap([]snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
		{Port: 80, Protocol: "tcp"},
	})

	checker := presencecheck.New(required, &buf)
	results := checker.Check(snap)
	checker.Report(results)

	output := buf.String()
	if !strings.Contains(output, "443") {
		t.Errorf("expected 443 to appear as missing, got: %s", output)
	}
	if strings.Contains(output, "22") {
		t.Errorf("did not expect port 22 in output, got: %s", output)
	}
	if strings.Contains(output, "80") {
		t.Errorf("did not expect port 80 in output, got: %s", output)
	}
}

// TestIntegration_EmptyRequired_NoOutput ensures that when no ports are
// required, no output is produced regardless of snapshot contents.
func TestIntegration_EmptyRequired_NoOutput(t *testing.T) {
	var buf bytes.Buffer

	snap := makeSnap([]snapshot.Entry{
		{Port: 22, Protocol: "tcp"},
	})
	checker := presencecheck.New(nil, &buf)
	results := checker.Check(snap)
	checker.Report(results)

	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty required list, got: %s", buf.String())
	}
}
