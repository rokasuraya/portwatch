package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// LogStage returns a StageFunc that writes a human-readable diff summary to w.
// If w is nil, os.Stdout is used.
func LogStage(w io.Writer) StageFunc {
	if w == nil {
		w = os.Stdout
	}
	return func(_ context.Context, opened, closed []snapshot.Entry) error {
		if len(opened) == 0 && len(closed) == 0 {
			return nil
		}
		ts := time.Now().UTC().Format(time.RFC3339)
		for _, e := range opened {
			fmt.Fprintf(w, "%s OPENED %s/%d\n", ts, e.Proto, e.Port)
		}
		for _, e := range closed {
			fmt.Fprintf(w, "%s CLOSED %s/%d\n", ts, e.Proto, e.Port)
		}
		return nil
	}
}

// NoopStage is a StageFunc that does nothing. Useful as a placeholder in
// tests or feature-flagged builds.
func NoopStage(_ context.Context, _, _ []snapshot.Entry) error { return nil }
