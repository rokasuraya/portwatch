// Package pipeline provides a composable scan-to-alert processing pipeline
// for portwatch.
//
// A Pipeline is constructed with a ScanFunc that produces port entries, a
// snapshot.Store for state persistence, and one or more StageFunc handlers
// that receive the opened/closed diff after each scan tick.
//
// Typical usage:
//
//	p := pipeline.New(myScanner, store,
//	    alertStage,
//	    metricsStage,
//	    auditStage,
//	)
//	duration, err := p.Tick(ctx)
package pipeline
