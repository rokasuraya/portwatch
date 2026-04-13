package labelmap_test

import (
	"testing"

	"github.com/user/portwatch/internal/labelmap"
)

func TestNew_BuiltInLabels(t *testing.T) {
	lm := labelmap.New(nil)

	cases := []struct {
		port  uint16
		want  string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{6379, "redis"},
	}
	for _, tc := range cases {
		got, ok := lm.Lookup(tc.port)
		if !ok {
			t.Errorf("port %d: expected found, got missing", tc.port)
		}
		if got != tc.want {
			t.Errorf("port %d: got %q, want %q", tc.port, got, tc.want)
		}
	}
}

func TestNew_ExtraOverridesBuiltIn(t *testing.T) {
	lm := labelmap.New(map[uint16]string{80: "my-http"})

	got, ok := lm.Lookup(80)
	if !ok {
		t.Fatal("expected port 80 to be found")
	}
	if got != "my-http" {
		t.Errorf("got %q, want %q", got, "my-http")
	}
}

func TestLabel_UnknownPort(t *testing.T) {
	lm := labelmap.New(nil)

	got := lm.Label(9999)
	want := "port/9999"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLabel_KnownPort(t *testing.T) {
	lm := labelmap.New(nil)

	got := lm.Label(22)
	if got != "ssh" {
		t.Errorf("got %q, want %q", got, "ssh")
	}
}

func TestRegister_AddsNewEntry(t *testing.T) {
	lm := labelmap.New(nil)
	lm.Register(9200, "elasticsearch")

	got, ok := lm.Lookup(9200)
	if !ok {
		t.Fatal("expected port 9200 to be found after Register")
	}
	if got != "elasticsearch" {
		t.Errorf("got %q, want %q", got, "elasticsearch")
	}
}

func TestLen_IncludesBuiltInAndExtra(t *testing.T) {
	extra := map[uint16]string{9200: "elasticsearch", 9300: "elasticsearch-transport"}
	lm := labelmap.New(extra)

	// builtIn has 19 entries; extra adds 2 new ones.
	if lm.Len() < 19 {
		t.Errorf("Len() = %d, want at least 19", lm.Len())
	}
}

func TestLookup_MissingReturnsFalse(t *testing.T) {
	lm := labelmap.New(nil)

	_, ok := lm.Lookup(65535)
	if ok {
		t.Error("expected ok=false for unregistered port 65535")
	}
}
