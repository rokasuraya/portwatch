package rotation

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestShouldRotate_NoFile(t *testing.T) {
	r := New("/tmp/portwatch_nonexistent_xyz.json", Options{MaxBytes: 100})
	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false for missing file")
	}
}

func TestShouldRotate_UnderLimit(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")
	writeFile(t, p, `{"small":true}`)

	r := New(p, Options{MaxBytes: 1024})
	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false when under size limit")
	}
}

func TestShouldRotate_OverLimit(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")
	writeFile(t, p, "x")

	r := New(p, Options{MaxBytes: 1})
	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true when at size limit")
	}
}

func TestShouldRotate_ByAge(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")
	writeFile(t, p, "data")

	r := New(p, Options{MaxAge: time.Millisecond})
	r.now = func() time.Time { return time.Now().Add(time.Hour) }

	ok, err := r.ShouldRotate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true when file is older than MaxAge")
	}
}

func TestRotate_RenamesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")
	writeFile(t, p, "original")

	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	r := New(p, Options{})
	r.now = func() time.Time { return fixed }

	if err := r.Rotate(); err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("original file should not exist after rotation")
	}
	expected := filepath.Join(dir, "state.20240601T120000Z.json")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("rotated file not found at %s: %v", expected, err)
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")

	// Create two existing backups.
	writeFile(t, filepath.Join(dir, "state.20240101T000000Z.json"), "old1")
	writeFile(t, filepath.Join(dir, "state.20240102T000000Z.json"), "old2")
	writeFile(t, p, "current")

	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	r := New(p, Options{MaxBackups: 2})
	r.now = func() time.Time { return fixed }

	if err := r.Rotate(); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "state.*.json"))
	if len(matches) != 2 {
		t.Fatalf("expected 2 backups after pruning, got %d: %v", len(matches), matches)
	}
}
