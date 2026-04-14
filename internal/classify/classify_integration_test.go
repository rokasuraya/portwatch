package classify_test

import (
	"testing"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/scanner"
)

// TestIntegration_RoundTrip verifies that applying a classifier over a mixed
// set of entries produces the expected distribution of risk levels.
func TestIntegration_RoundTrip(t *testing.T) {
	rules := []classify.Rule{
		{Port: 5432, Protocol: "tcp", Level: classify.LevelHigh}, // postgres flagged as high
	}
	c := classify.New(rules)

	entries := []scanner.Entry{
		{Port: 22, Protocol: "tcp"},   // built-in high
		{Port: 443, Protocol: "tcp"},  // medium (privileged)
		{Port: 8443, Protocol: "tcp"}, // low (ephemeral)
		{Port: 5432}, // overridden to high
		{Port: 5432, Protocol: "udp"}, // no override → low
	}

	result := c.Apply(entries)

	want := map[scanner.Entry]classify.Level{
		{Port: 22, Protocol: "tcp"}:   classify.LevelHigh,
		{Port: 443, Protocol: "tcp"}:  classify.LevelMedium,
		{Port: 8443, Protocol: "tcp"}: classify.LevelLow,
		{Port: 5432, Protocol: "tcp"}: classify.LevelHigh,
		{Port: 5432, Protocol: "udp"}: classify.LevelLow,
	}

	for e, wantLvl := range want {
		if got := result[e]; got != wantLvl {
			t.Errorf("entry %v: want %s, got %s", e, wantLvl, got)
		}
	}
}

// TestIntegration_EmptyEntries ensures Apply handles an empty slice gracefully.
func TestIntegration_EmptyEntries(t *testing.T) {
	c := classify.New(nil)
	result := c.Apply([]scanner.Entry{})
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(result))
	}
}
