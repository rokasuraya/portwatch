package classify_test

import (
	"testing"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/scanner"
)

func entry(port int, proto string) scanner.Entry {
	return scanner.Entry{Port: port, Protocol: proto}
}

func TestNew_NoRules(t *testing.T) {
	c := classify.New(nil)
	if c == nil {
		t.Fatal("expected non-nil classifier")
	}
}

func TestClassify_SSHIsHigh(t *testing.T) {
	c := classify.New(nil)
	if got := c.Classify(entry(22, "tcp")); got != classify.LevelHigh {
		t.Fatalf("port 22: want high, got %s", got)
	}
}

func TestClassify_RDPIsHigh(t *testing.T) {
	c := classify.New(nil)
	if got := c.Classify(entry(3389, "tcp")); got != classify.LevelHigh {
		t.Fatalf("port 3389: want high, got %s", got)
	}
}

func TestClassify_PrivilegedIsMedium(t *testing.T) {
	c := classify.New(nil)
	// port 80 is <1024 but not in the high list
	if got := c.Classify(entry(80, "tcp")); got != classify.LevelMedium {
		t.Fatalf("port 80: want medium, got %s", got)
	}
}

func TestClassify_EphemeralIsLow(t *testing.T) {
	c := classify.New(nil)
	if got := c.Classify(entry(8080, "tcp")); got != classify.LevelLow {
		t.Fatalf("port 8080: want low, got %s", got)
	}
}

func TestClassify_RuleOverridesDefault(t *testing.T) {
	rules := []classify.Rule{
		{Port: 8080, Protocol: "tcp", Level: classify.LevelHigh},
	}
	c := classify.New(rules)
	if got := c.Classify(entry(8080, "tcp")); got != classify.LevelHigh {
		t.Fatalf("override: want high, got %s", got)
	}
}

func TestClassify_RuleProtocolDistinct(t *testing.T) {
	rules := []classify.Rule{
		{Port: 8080, Protocol: "tcp", Level: classify.LevelHigh},
	}
	c := classify.New(rules)
	// udp:8080 has no override — falls back to low
	if got := c.Classify(entry(8080, "udp")); got != classify.LevelLow {
		t.Fatalf("protocol mismatch: want low, got %s", got)
	}
}

func TestApply_ReturnsMapForAll(t *testing.T) {
	c := classify.New(nil)
	entries := []scanner.Entry{
		entry(22, "tcp"),
		entry(80, "tcp"),
		entry(9000, "tcp"),
	}
	result := c.Apply(entries)
	if len(result) != len(entries) {
		t.Fatalf("want %d entries, got %d", len(entries), len(result))
	}
	if result[entry(22, "tcp")] != classify.LevelHigh {
		t.Error("port 22 should be high")
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		lvl  classify.Level
		want string
	}{
		{classify.LevelLow, "low"},
		{classify.LevelMedium, "medium"},
		{classify.LevelHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.lvl.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.lvl, got, tc.want)
		}
	}
}
