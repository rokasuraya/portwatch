package acknowledge_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/acknowledge"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ack.json")
}

func TestNew_NoExistingFile(t *testing.T) {
	a, err := acknowledge.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(a.All()); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestAck_IsAcked_BasicFlow(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	if a.IsAcked("tcp", 22) {
		t.Fatal("port should not be acked before Ack call")
	}
	if err := a.Ack("tcp", 22, time.Time{}, "planned maintenance"); err != nil {
		t.Fatalf("Ack returned error: %v", err)
	}
	if !a.IsAcked("tcp", 22) {
		t.Fatal("port should be acked after Ack call")
	}
}

func TestAck_ProtocolDistinct(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	_ = a.Ack("tcp", 53, time.Time{}, "")
	if a.IsAcked("udp", 53) {
		t.Fatal("udp:53 should not be acked when only tcp:53 was acked")
	}
}

func TestAck_Expiry(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	past := time.Now().Add(-time.Second)
	_ = a.Ack("tcp", 80, past, "expired")
	if a.IsAcked("tcp", 80) {
		t.Fatal("expired acknowledgement should not be considered active")
	}
}

func TestAck_NoExpiry(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	_ = a.Ack("tcp", 443, time.Time{}, "permanent")
	if !a.IsAcked("tcp", 443) {
		t.Fatal("zero-expiry acknowledgement should always be active")
	}
}

func TestRemove_ClearsAck(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	_ = a.Ack("tcp", 8080, time.Time{}, "")
	if err := a.Remove("tcp", 8080); err != nil {
		t.Fatalf("Remove returned error: %v", err)
	}
	if a.IsAcked("tcp", 8080) {
		t.Fatal("port should not be acked after Remove")
	}
}

func TestRemove_NonexistentEntry(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	if err := a.Remove("tcp", 9999); err != nil {
		t.Fatalf("Remove of non-existent entry should not error, got: %v", err)
	}
}

func TestAck_PersistsAcrossSessions(t *testing.T) {
	path := tempPath(t)
	a1, _ := acknowledge.New(path)
	_ = a1.Ack("tcp", 9090, time.Time{}, "persisted")

	a2, err := acknowledge.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if !a2.IsAcked("tcp", 9090) {
		t.Fatal("acknowledgement should survive a reload from disk")
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := acknowledge.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON file")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	a, _ := acknowledge.New(tempPath(t))
	_ = a.Ack("tcp", 22, time.Time{}, "a")
	_ = a.Ack("udp", 53, time.Time{}, "b")
	if got := len(a.All()); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}
