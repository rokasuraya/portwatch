// Package portpulse provides a rolling-window rate tracker for port-change
// events observed during portwatch scans.
//
// A Tracker accepts snapshot diffs, counts opened and closed port events,
// and prunes entries that fall outside a configurable time window.  The
// Rate method returns the total number of change events still within the
// window, giving callers a lightweight "pulse" of how volatile the
// monitored host's port surface has been recently.
//
// Usage:
//
//	tr := portpulse.New(5*time.Minute, os.Stdout)
//	tr.Observe(diff)
//	fmt.Println(tr.Rate()) // events in last 5 minutes
//	tr.Report()            // human-readable summary
package portpulse
