package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestNewStore_NoExistingFile(t *testing.T) {
	st, err := snapshot.NewStore(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestStore_SetAndCurrent(t *testing.T) {
	path := tempPath(t)
	st, _ := snapshot.NewStore(path)

	snap := snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 80}})
	if err := st.Set(snap); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	got := st.Current()
	if len(got.Entries) != 1 || got.Entries[0].Port != 80 {
		t.Fatalf("unexpected current snapshot: %+v", got)
	}
}

func TestStore_SetPromotesPrevious(t *testing.T) {
	path := tempPath(t)
	st, _ := snapshot.NewStore(path)

	first := snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 22}})
	_ = st.Set(first)

	second := snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 443}})
	_ = st.Set(second)

	prev := st.Previous()
	if len(prev.Entries) != 1 || prev.Entries[0].Port != 22 {
		t.Fatalf("expected previous to hold port 22, got %+v", prev)
	}
}

func TestStore_PersistsToDisk(t *testing.T) {
	path := tempPath(t)
	st, _ := snapshot.NewStore(path)
	snap := snapshot.New([]snapshot.Entry{{Protocol: "udp", Port: 53}})
	_ = st.Set(snap)

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestNewStore_LoadsExistingFile(t *testing.T) {
	path := tempPath(t)
	st1, _ := snapshot.NewStore(path)
	_ = st1.Set(snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 9090}}))

	st2, err := snapshot.NewStore(path)
	if err != nil {
		t.Fatalf("unexpected error on reload: %v", err)
	}
	if len(st2.Current().Entries) != 1 || st2.Current().Entries[0].Port != 9090 {
		t.Fatalf("expected reloaded port 9090, got %+v", st2.Current())
	}
}
