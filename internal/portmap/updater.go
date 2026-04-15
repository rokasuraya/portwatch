package portmap

import "github.com/user/portwatch/internal/snapshot"

// Updater applies snapshot diffs to a PortMap, keeping it in sync with
// the latest scan results.
type Updater struct {
	pm     *PortMap
	labels func(port int, proto string) string
}

// NewUpdater returns an Updater that writes changes into pm.
// The labels function is called to resolve a human-readable name for each
// port; pass nil to use an empty label.
func NewUpdater(pm *PortMap, labels func(int, string) string) *Updater {
	if labels == nil {
		labels = func(int, string) string { return "" }
	}
	return &Updater{pm: pm, labels: labels}
}

// Apply processes opened and closed entries from a snapshot diff, updating
// the underlying PortMap accordingly.
func (u *Updater) Apply(opened, closed []snapshot.Entry) {
	for _, e := range opened {
		u.pm.Set(e.Port, e.Protocol, u.labels(e.Port, e.Protocol), true)
	}
	for _, e := range closed {
		u.pm.Delete(e.Port, e.Protocol)
	}
}
