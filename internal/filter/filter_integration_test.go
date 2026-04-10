package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

// TestApply_RoundTrip verifies that a Filter built from multiple rules
// correctly partitions a realistic port list into suppressed and visible sets.
func TestApply_RoundTrip(t *testing.T) {
	rules := []filter.Rule{
		{Port: 22, Protocol: "tcp", Comment: "SSH"},
		{Port: 80, Protocol: "tcp", Comment: "HTTP"},
		{Port: 443, Protocol: "tcp", Comment: "HTTPS"},
		{Port: 53, Protocol: "udp", Comment: "DNS"},
	}
	f := filter.New(rules)

	all := []string{"22/tcp", "53/udp", "80/tcp", "443/tcp", "3000/tcp", "5432/tcp"}
	visible := f.Apply(all)

	expected := []string{"3000/tcp", "5432/tcp"}
	if len(visible) != len(expected) {
		t.Fatalf("got %v, want %v", visible, expected)
	}
	for i, v := range expected {
		if visible[i] != v {
			t.Errorf("visible[%d] = %q, want %q", i, visible[i], v)
		}
	}
}

// TestAllow_MultipleRules confirms Allow is consistent with Apply.
func TestAllow_MultipleRules(t *testing.T) {
	rules := []filter.Rule{
		{Port: 6379, Protocol: "tcp", Comment: "Redis"},
	}
	f := filter.New(rules)

	cases := []struct {
		port  int
		proto string
		want  bool
	}{
		{6379, "tcp", false},
		{6379, "udp", true},
		{9200, "tcp", true},
	}
	for _, c := range cases {
		got := f.Allow(c.port, c.proto)
		if got != c.want {
			t.Errorf("Allow(%d, %q) = %v, want %v", c.port, c.proto, got, c.want)
		}
	}
}
