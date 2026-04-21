package portdiff_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portdiff"
)

func TestFormat_NoChanges(t *testing.T) {
	var sb strings.Builder
	err := portdiff.Format(&sb, portdiff.Diff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no port changes") {
		t.Errorf("expected no-change message, got: %s", sb.String())
	}
}

func TestFormat_Opened(t *testing.T) {
	d := portdiff.Diff{
		Opened: []portdiff.Entry{{Op: "opened", Port: 80, Protocol: "tcp", Label: "http"}},
	}
	var sb strings.Builder
	if err := portdiff.Format(&sb, d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := sb.String()
	if !strings.HasPrefix(got, "+ ") {
		t.Errorf("expected line to start with '+ ', got: %s", got)
	}
	if !strings.Contains(got, "http") {
		t.Errorf("expected label in output, got: %s", got)
	}
}

func TestFormat_Closed(t *testing.T) {
	d := portdiff.Diff{
		Closed: []portdiff.Entry{{Op: "closed", Port: 443, Protocol: "tcp"}},
	}
	var sb strings.Builder
	if err := portdiff.Format(&sb, d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "- closed tcp/443") {
		t.Errorf("unexpected output: %s", sb.String())
	}
}

func TestFormatJSON_Structure(t *testing.T) {
	d := portdiff.Diff{
		Opened: []portdiff.Entry{{Op: "opened", Port: 22, Protocol: "tcp", Label: "ssh"}},
		Closed: []portdiff.Entry{{Op: "closed", Port: 8080, Protocol: "tcp", Label: ""}},
	}
	var sb strings.Builder
	if err := portdiff.FormatJSON(&sb, d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := sb.String()
	if !strings.Contains(got, `"opened":[`) {
		t.Errorf("missing opened array: %s", got)
	}
	if !strings.Contains(got, `"closed":[`) {
		t.Errorf("missing closed array: %s", got)
	}
	if !strings.Contains(got, `"label":"ssh"`) {
		t.Errorf("missing label field: %s", got)
	}
}

func TestFormatJSON_Empty(t *testing.T) {
	var sb strings.Builder
	if err := portdiff.FormatJSON(&sb, portdiff.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(sb.String())
	if got != `{"opened":[],"closed":[]}` {
		t.Errorf("unexpected JSON: %s", got)
	}
}
