package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(proto string, port int) snapshot.Entry {
	return snapshot.Entry{Protocol: proto, Port: port}
}

func TestCollect_RecordsOpened(t *testing.T) {
	var buf bytes.Buffer
	c := audit.NewCollector(audit.New(&buf))
	c.Collect([]snapshot.Entry{makeEntry("tcp", 9000)}, nil)

	var e audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e)
	if e.Event != "opened" || e.Port != 9000 {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestCollect_RecordsClosed(t *testing.T) {
	var buf bytes.Buffer
	c := audit.NewCollector(audit.New(&buf))
	c.Collect(nil, []snapshot.Entry{makeEntry("udp", 161)})

	var e audit.Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e)
	if e.Event != "closed" || e.Port != 161 {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestCollect_MixedDiff(t *testing.T) {
	var buf bytes.Buffer
	c := audit.NewCollector(audit.New(&buf))
	c.Collect(
		[]snapshot.Entry{makeEntry("tcp", 8080)},
		[]snapshot.Entry{makeEntry("tcp", 22)},
	)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestCollect_EmptyDiff(t *testing.T) {
	var buf bytes.Buffer
	c := audit.NewCollector(audit.New(&buf))
	c.Collect(nil, nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got %q", buf.String())
	}
}
