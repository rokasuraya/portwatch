package reporter

import (
	"testing"
	"time"
)

func TestBuildReport_PopulatesFields(t *testing.T) {
	opened := []int{8080, 9090}
	closed := []int{443}
	current := []int{22, 80, 8080, 9090}

	before := time.Now().UTC()
	report := BuildReport(opened, closed, current)
	after := time.Now().UTC()

	if report.Timestamp.Before(before) || report.Timestamp.After(after) {
		t.Error("timestamp outside expected range")
	}
	if len(report.OpenedPorts) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(report.OpenedPorts))
	}
	if len(report.ClosedPorts) != 1 {
		t.Errorf("expected 1 closed port, got %d", len(report.ClosedPorts))
	}
	if report.TotalOpen != 4 {
		t.Errorf("expected TotalOpen=4, got %d", report.TotalOpen)
	}
}

func TestBuildReport_EmptySlices(t *testing.T) {
	report := BuildReport(nil, nil, []int{22})
	if report.OpenedPorts == nil {
		t.Error("expected non-nil OpenedPorts slice")
	}
	if report.ClosedPorts == nil {
		t.Error("expected non-nil ClosedPorts slice")
	}
	if report.TotalOpen != 1 {
		t.Errorf("expected TotalOpen=1, got %d", report.TotalOpen)
	}
}

func TestBuildReport_IsolatesSlices(t *testing.T) {
	opened := []int{8080}
	report := BuildReport(opened, nil, nil)
	opened[0] = 9999
	if report.OpenedPorts[0] == 9999 {
		t.Error("report should not share backing array with input slice")
	}
}
