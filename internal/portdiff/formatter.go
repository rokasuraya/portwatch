package portdiff

import (
	"fmt"
	"io"
	"strings"
)

// Format writes a multi-line human-readable representation of d to w.
// Each entry is prefixed with "+" for opened and "-" for closed.
func Format(w io.Writer, d Diff) error {
	if d.IsEmpty() {
		_, err := fmt.Fprintln(w, "no port changes detected")
		return err
	}

	var sb strings.Builder
	for _, e := range d.Opened {
		sb.WriteString("+ ")
		sb.WriteString(e.String())
		sb.WriteByte('\n')
	}
	for _, e := range d.Closed {
		sb.WriteString("- ")
		sb.WriteString(e.String())
		sb.WriteByte('\n')
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// FormatJSON writes a JSON representation of d to w.
// The output is a single JSON object with "opened" and "closed" arrays.
func FormatJSON(w io.Writer, d Diff) error {
	var sb strings.Builder
	sb.WriteString(`{"opened":[`)
	for i, e := range d.Opened {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(entryJSON(e))
	}
	sb.WriteString(`],"closed":[`)
	for i, e := range d.Closed {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(entryJSON(e))
	}
	sb.WriteString("]}")  
	_, err := fmt.Fprintln(w, sb.String())
	return err
}

func entryJSON(e Entry) string {
	return fmt.Sprintf(`{"op":%q,"port":%d,"protocol":%q,"label":%q}`,
		e.Op, e.Port, e.Protocol, e.Label)
}
