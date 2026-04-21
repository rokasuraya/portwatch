// Package portrank provides a composable scoring system for ranking open ports
// by risk. Callers supply one or more ScorerFunc implementations that each
// contribute a partial float64 score and an optional human-readable reason.
//
// Built-in scorers:
//   - PrivilegedPortScorer  – penalises ports below 1024
//   - WellKnownRiskyPortScorer – penalises commonly exploited services
//   - ProtocolScorer        – adds a small penalty for UDP
//
// Results are returned sorted by descending composite score so the highest-risk
// ports appear first.
package portrank
