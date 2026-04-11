// Package healthcheck provides a lightweight HTTP liveness endpoint for the
// portwatch daemon.
//
// Usage:
//
//	srv := healthcheck.New(":9090")
//	go srv.ListenAndServe()
//	// ... after each scan cycle:
//	srv.RecordScan()
//	// on shutdown:
//	srv.Shutdown()
//
// GET /healthz returns a JSON payload:
//
//	{"ok":true,"uptime":"2m30s","last_scan":"...","scan_count":15}
package healthcheck
