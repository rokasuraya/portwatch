package history_test

import (
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/state"
)

func TestCollect_RecordsOpened(t *testing.T) {
	h, _ := history.New(tempPath(t), 20)
	c := history.NewCollector(h)

	diff := state.Diff{
		Opened: []state.Entry{{Proto: "tcp", Port: 80}},
	}
	if err := c.Collect(diff); err != nil {
		t.Fatalf("Collect: %v", err)
	}
	events := h.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != "opened" || events[0].Port != 80 {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestCollect_RecordsClosed(t *testing.T) {
	h, _ := history.New(tempPath(t), 20)
	c := history.NewCollector(h)

	diff := state.Diff{
		Closed: []state.Entry{{Proto: "udp", Port: 53}},
	}
	_ = c.Collect(diff)
	events := h.Events()
	if len(events) != 1 || events[0].Kind != "closed" {
		t.Errorf("expected 1 closed event, got %+v", events)
	}
}

func TestCollect_MixedDiff(t *testing.T) {
	h, _ := history.New(tempPath(t), 20)
	c := history.NewCollector(h)

	diff := state.Diff{
		Opened: []state.Entry{{Proto: "tcp", Port: 443}, {Proto: "tcp", Port: 8080}},
		Closed: []state.Entry{{Proto: "tcp", Port: 22}},
	}
	_ = c.Collect(diff)
	if got := len(h.Events()); got != 3 {
		t.Fatalf("expected 3 events, got %d", got)
	}
}

func TestCollect_EmptyDiff(t *testing.T) {
	h, _ := history.New(tempPath(t), 20)
	c := history.NewCollector(h)

	_ = c.Collect(state.Diff{})
	if got := len(h.Events()); got != 0 {
		t.Fatalf("expected 0 events, got %d", got)
	}
}
