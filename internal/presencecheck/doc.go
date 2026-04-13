// Package presencecheck provides a Checker that validates a set of required
// ports are open in a given snapshot. It is useful for ensuring critical
// services (e.g. SSH on port 22, HTTPS on port 443) remain reachable.
//
// Usage:
//
//	required := []snapshot.Entry{
//		{Port: 22, Protocol: "tcp"},
//		{Port: 443, Protocol: "tcp"},
//	}
//	checker := presencecheck.New(required, nil)
//	results := checker.Check(snap)
//	checker.Report(results)
package presencecheck
