package portquota_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portquota"
)

func TestIntegration_BreachThenRecover(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(3, &buf)

	// Gradually grow past limit.
	for i, ports := range [][]int{
		{80},
		{80, 443},
		{80, 443, 8080},
		{80, 443, 8080, 9090}, // breach
	} {
		breach := q.Check(makeSnap(ports))
		if i < 3 && breach {
			t.Errorf("step %d: unexpected breach", i)
		}
		if i == 3 && !breach {
			t.Errorf("step %d: expected breach", i)
		}
	}

	if !strings.Contains(buf.String(), "exceeds limit 3") {
		t.Errorf("expected limit in warning, got: %s", buf.String())
	}

	// Recover.
	if q.Check(makeSnap([]int{80})) {
		t.Error("expected no breach after recovery")
	}

	// Breach again — warning must reappear.
	buf.Reset()
	if !q.Check(makeSnap([]int{80, 443, 8080, 9090})) {
		t.Error("expected breach on second overage")
	}
	if !strings.Contains(buf.String(), "exceeds limit") {
		t.Error("expected warning after second breach")
	}
}

func TestIntegration_DynamicMaxAdjustment(t *testing.T) {
	var buf bytes.Buffer
	q := portquota.New(10, &buf)
	snap := makeSnap([]int{80, 443, 8080, 9090, 3000})

	if q.Check(snap) {
		t.Error("expected no breach with generous limit")
	}

	q.SetMax(3)
	if !q.Check(snap) {
		t.Error("expected breach after tightening limit")
	}
}
