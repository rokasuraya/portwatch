package portlabel

import (
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port uint16, proto, label string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto, Label: label}
}

func TestAnnotate_SetsLabelForKnownPort(t *testing.T) {
	a := NewAnnotator(New(nil))
	entries := []snapshot.Entry{makeEntry(22, "tcp", "")}
	out := a.Annotate(entries)
	if out[0].Label != "SSH" {
		t.Fatalf("expected SSH, got %q", out[0].Label)
	}
}

func TestAnnotate_SetsGenericLabelForUnknownPort(t *testing.T) {
	a := NewAnnotator(New(nil))
	entries := []snapshot.Entry{makeEntry(9999, "tcp", "")}
	out := a.Annotate(entries)
	if out[0].Label == "" {
		t.Fatal("expected a non-empty label for unknown port")
	}
}

func TestAnnotate_PreservesExistingLabel(t *testing.T) {
	a := NewAnnotator(New(nil))
	entries := []snapshot.Entry{makeEntry(22, "tcp", "AlreadySet")}
	out := a.Annotate(entries)
	if out[0].Label != "AlreadySet" {
		t.Fatalf("expected AlreadySet, got %q", out[0].Label)
	}
}

func TestAnnotate_DoesNotMutateInput(t *testing.T) {
	a := NewAnnotator(New(nil))
	orig := []snapshot.Entry{makeEntry(80, "tcp", "")}
	_ = a.Annotate(orig)
	if orig[0].Label != "" {
		t.Fatal("input slice was mutated")
	}
}

func TestAnnotate_EmptySlice(t *testing.T) {
	a := NewAnnotator(New(nil))
	out := a.Annotate(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d entries", len(out))
	}
}

func TestNewAnnotator_NilLabeler(t *testing.T) {
	a := NewAnnotator(nil)
	if a == nil {
		t.Fatal("expected non-nil Annotator")
	}
	entries := []snapshot.Entry{makeEntry(443, "tcp", "")}
	out := a.Annotate(entries)
	if out[0].Label != "HTTPS" {
		t.Fatalf("expected HTTPS, got %q", out[0].Label)
	}
}
