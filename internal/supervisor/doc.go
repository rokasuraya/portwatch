// Package supervisor wires together the core portwatch components and
// drives a single scan tick.
//
// It sits between the daemon (which owns the scheduling loop) and the
// lower-level packages (scanner, state, metrics). Keeping this
// orchestration in its own package makes the daemon loop trivial and
// keeps each component independently testable.
//
// Typical usage:
//
//	sv := supervisor.New(cfg, supervisor.Components{
//		Scanner:  sc,
//		State:    st,
//		Metrics:  m,
//		OnChange: alertFn,
//	}, logger)
//
//	if err := sv.Tick(ctx); err != nil {
//		log.Println("tick error:", err)
//	}
package supervisor
