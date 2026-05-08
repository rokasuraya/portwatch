package portlabel

import "github.com/user/portwatch/internal/snapshot"

// Annotator enriches snapshot entries with label metadata.
type Annotator struct {
	labeler *Labeler
}

// NewAnnotator returns an Annotator backed by the given Labeler.
func NewAnnotator(l *Labeler) *Annotator {
	if l == nil {
		l = New(nil)
	}
	return &Annotator{labeler: l}
}

// Annotate returns a copy of entries with the Label field populated.
// Entries that already carry a non-empty label are left unchanged.
func (a *Annotator) Annotate(entries []snapshot.Entry) []snapshot.Entry {
	out := make([]snapshot.Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Label != "" {
			continue
		}
		lbl, _ := a.labeler.Lookup(e.Port, e.Protocol)
		out[i].Label = lbl.Name
	}
	return out
}
