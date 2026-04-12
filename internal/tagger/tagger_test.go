package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func TestNew_NilExtra(t *testing.T) {
	tg := tagger.New(nil)
	if tg == nil {
		t.Fatal("expected non-nil Tagger")
	}
}

func TestLabel_WellKnownPort(t *testing.T) {
	tg := tagger.New(nil)
	got := tg.Label(22)
	if got != "ssh" {
		t.Errorf("expected ssh, got %q", got)
	}
}

func TestLabel_UnknownPort(t *testing.T) {
	tg := tagger.New(nil)
	got := tg.Label(9999)
	if got != "port/9999" {
		t.Errorf("expected port/9999, got %q", got)
	}
}

func TestLabel_ExtraOverridesBuiltIn(t *testing.T) {
	extra := map[int]string{80: "my-app"}
	tg := tagger.New(extra)
	got := tg.Label(80)
	if got != "my-app" {
		t.Errorf("expected my-app, got %q", got)
	}
}

func TestLabel_ExtraOnlyPort(t *testing.T) {
	extra := map[int]string{12345: "custom-svc"}
	tg := tagger.New(extra)
	got := tg.Label(12345)
	if got != "custom-svc" {
		t.Errorf("expected custom-svc, got %q", got)
	}
}

func TestTag_SetsLabelOnEntries(t *testing.T) {
	tg := tagger.New(nil)
	entries := []scanner.Entry{
		{Port: 443, Proto: "tcp"},
		{Port: 3306, Proto: "tcp"},
		{Port: 7777, Proto: "tcp"},
	}
	tg.Tag(entries)

	cases := []struct {
		idx  int
		want string
	}{
		{0, "https"},
		{1, "mysql"},
		{2, "port/7777"},
	}
	for _, c := range cases {
		if entries[c.idx].Label != c.want {
			t.Errorf("entry[%d].Label = %q, want %q", c.idx, entries[c.idx].Label, c.want)
		}
	}
}

func TestTag_EmptySlice(t *testing.T) {
	tg := tagger.New(nil)
	var entries []scanner.Entry
	tg.Tag(entries) // must not panic
}
