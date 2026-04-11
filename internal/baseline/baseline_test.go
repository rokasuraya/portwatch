package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_NoExistingFile(t *testing.T) {
	b, err := baseline.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(b.Entries))
	}
}

func TestApprove_PersistsEntries(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path)

	snap := makeSnap([]snapshot.Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	})
	if err := b.Approve(snap); err != nil {
		t.Fatalf("Approve error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
	if len(b.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(b.Entries))
	}
}

func TestApprove_SetsUpdatedTime(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	b, _ := baseline.New(tempPath(t))
	_ = b.Approve(makeSnap([]snapshot.Entry{{Port: 22, Protocol: "tcp"}}))
	if !b.Updated.After(before) {
		t.Errorf("Updated not set: %v", b.Updated)
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path)
	_ = b.Approve(makeSnap([]snapshot.Entry{{Port: 8080, Protocol: "tcp"}}))

	b2, err := baseline.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if len(b2.Entries) != 1 || b2.Entries[0].Port != 8080 {
		t.Errorf("unexpected entries after reload: %+v", b2.Entries)
	}
}

func TestDiff_DetectsUnexpected(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Approve(makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}}))

	snap := makeSnap([]snapshot.Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 9999, Protocol: "tcp"},
	})
	unexpected, missing := b.Diff(snap)
	if len(unexpected) != 1 || unexpected[0].Port != 9999 {
		t.Errorf("unexpected diff: %+v", unexpected)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %+v", missing)
	}
}

func TestDiff_DetectsMissing(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Approve(makeSnap([]snapshot.Entry{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}))

	snap := makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}})
	unexpected, missing := b.Diff(snap)
	if len(unexpected) != 0 {
		t.Errorf("expected no unexpected, got %+v", unexpected)
	}
	if len(missing) != 1 || missing[0].Port != 443 {
		t.Errorf("missing diff wrong: %+v", missing)
	}
}

func TestDiff_EmptyBaseline_AllUnexpected(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	snap := makeSnap([]snapshot.Entry{{Port: 22, Protocol: "tcp"}})
	unexpected, missing := b.Diff(snap)
	if len(unexpected) != 1 {
		t.Errorf("expected 1 unexpected, got %d", len(unexpected))
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing, got %d", len(missing))
	}
}
