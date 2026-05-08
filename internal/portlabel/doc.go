// Package portlabel provides a registry that maps port+protocol pairs to
// human-readable labels and categories.
//
// A Labeler is initialised with a set of well-known built-in entries
// (SSH, HTTP, HTTPS, …) and can be extended at runtime via Register or
// by supplying an overrides map to New.
//
// The companion Annotator type enriches snapshot.Entry slices in-place,
// filling the Label field for any entry whose label is currently empty.
//
// Example:
//
//	labeler := portlabel.New(nil)
//	annotator := portlabel.NewAnnotator(labeler)
//	annotated := annotator.Annotate(entries)
package portlabel
