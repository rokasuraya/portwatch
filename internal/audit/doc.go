// Package audit provides a structured, append-only audit trail for portwatch.
//
// Each port change event (opened or closed) is serialised as a JSON line and
// written to a configurable io.Writer — typically a dedicated audit log file.
//
// Usage:
//
//	f, _ := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
//	a := audit.New(f)
//	c := audit.NewCollector(a)
//	c.Collect(opened, closed)
package audit
