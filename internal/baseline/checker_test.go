package baseline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/snapshot"
)

func TestChecker_CleanResult(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Approve(makeSnap([]snapshot.Entry{
		{Port: 80, Protocol: "tcp"},
	}))

	var buf bytes.Buffer
	c := baseline.NewChecker(b, &buf)
	res := c.Check(makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}}))

	if !res.Clean {
		t.Errorf("expected clean result")
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK in output, got: %s", buf.String())
	}
}

func TestChecker_UnexpectedPort(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Approve(makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}}))

	var buf bytes.Buffer
	c := baseline.NewChecker(b, &buf)
	res := c.Check(makeSnap([]snapshot.Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 4444, Protocol: "tcp"},
	}))

	if res.Clean {
		t.Errorf("expected deviation")
	}
	if len(res.Unexpected) != 1 || res.Unexpected[0].Port != 4444 {
		t.Errorf("wrong unexpected: %+v", res.Unexpected)
	}
	if !strings.Contains(buf.String(), "UNEXPECTED") {
		t.Errorf("expected UNEXPECTED in output, got: %s", buf.String())
	}
}

func TestChecker_MissingPort(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Approve(makeSnap([]snapshot.Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}))

	var buf bytes.Buffer
	c := baseline.NewChecker(b, &buf)
	res := c.Check(makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}}))

	if res.Clean {
		t.Errorf("expected deviation")
	}
	if len(res.Missing) != 1 || res.Missing[0].Port != 443 {
		t.Errorf("wrong missing: %+v", res.Missing)
	}
	if !strings.Contains(buf.String(), "MISSING") {
		t.Errorf("expected MISSING in output, got: %s", buf.String())
	}
}

func TestNewChecker_DefaultsToStdout(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	c := baseline.NewChecker(b, nil)
	if c == nil {
		t.Fatal("expected non-nil checker")
	}
}

func TestChecker_TimestampSet(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	var buf bytes.Buffer
	c := baseline.NewChecker(b, &buf)
	res := c.Check(makeSnap(nil))
	if res.Timestamp.IsZero() {
		t.Errorf("expected non-zero timestamp")
	}
}
