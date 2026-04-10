package reporter

import "time"

// BuildReport constructs a Report from opened/closed port slices and the
// current total number of open ports.
func BuildReport(opened, closed, currentOpen []int) Report {
	return Report{
		Timestamp:   time.Now().UTC(),
		OpenedPorts: copySlice(opened),
		ClosedPorts: copySlice(closed),
		TotalOpen:   len(currentOpen),
	}
}

func copySlice(s []int) []int {
	if len(s) == 0 {
		return []int{}
	}
	out := make([]int, len(s))
	copy(out, s)
	return out
}
