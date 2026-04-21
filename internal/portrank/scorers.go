package portrank

import (
	"fmt"

	"portwatch/internal/snapshot"
)

// PrivilegedPortScorer returns a higher score for ports below 1024.
func PrivilegedPortScorer() ScorerFunc {
	return func(e snapshot.Entry) (float64, string) {
		if e.Port < 1024 {
			return 2.0, fmt.Sprintf("privileged port %d", e.Port)
		}
		return 0, ""
	}
}

// WellKnownRiskyPortScorer bumps the score for commonly exploited ports.
func WellKnownRiskyPortScorer() ScorerFunc {
	risky := map[int]string{
		22:   "SSH",
		23:   "Telnet",
		3389: "RDP",
		445:  "SMB",
		1433: "MSSQL",
		3306: "MySQL",
	}
	return func(e snapshot.Entry) (float64, string) {
		if label, ok := risky[e.Port]; ok {
			return 3.0, fmt.Sprintf("well-known risky service: %s", label)
		}
		return 0, ""
	}
}

// ProtocolScorer adds a small penalty for UDP ports, which are harder to audit.
func ProtocolScorer() ScorerFunc {
	return func(e snapshot.Entry) (float64, string) {
		if e.Protocol == "udp" {
			return 0.5, "UDP protocol (harder to audit)"
		}
		return 0, ""
	}
}
