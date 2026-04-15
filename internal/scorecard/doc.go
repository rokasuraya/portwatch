// Package scorecard provides a lightweight risk-scoring layer for port diff
// events. It combines the number of opened/closed ports with their severity
// classification (via the classify package) to produce a numeric total and a
// human-readable band label (none / low / medium / high).
//
// Typical usage:
//
//	sc := scorecard.New(classifier)
//	score := sc.Evaluate(opened, closed)
//	fmt.Println(score.Label, score.Total)
package scorecard
