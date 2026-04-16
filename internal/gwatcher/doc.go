// Package gwatcher tracks port group membership across successive snapshots
// and emits join/leave events when a port moves into or out of a named group
// defined in the portgroup registry.
//
// Typical usage:
//
//	reg := portgroup.New()
//	reg.Define("web", []portgroup.Entry{{Port: 80, Proto: "tcp"}, {Port: 443, Proto: "tcp"}})
//	matcher := portgroup.NewMatcher(reg)
//	w := gwatcher.New(matcher, os.Stdout)
//	events := w.Observe(snap)
package gwatcher
