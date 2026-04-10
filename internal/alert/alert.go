package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      int
}

// Notifier sends alerts about port changes.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify evaluates a state.Diff and emits alerts for opened and closed ports.
func (n *Notifier) Notify(diff state.Diff) []Alert {
	var alerts []Alert
	now := time.Now()

	for _, port := range diff.Opened {
		a := Alert{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port %d was opened", port),
			Port:      port,
		}
		alerts = append(alerts, a)
		n.write(a)
	}

	for _, port := range diff.Closed {
		a := Alert{
			Timestamp: now,
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port %d was closed", port),
			Port:      port,
		}
		alerts = append(alerts, a)
		n.write(a)
	}

	return alerts
}

func (n *Notifier) write(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
	)
}
