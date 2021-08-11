package matrix

import "strings"

type TagName string

// HasNamespace returns whether the specified tagname has the namespace.
func (t TagName) HasNamespace(n string) bool {
	return strings.HasPrefix(string(t), n+".")
}

// Tags defined in the matrix namespace.
var (
	TagFavourite    TagName = "m.favourite"
	TagLowPriority  TagName = "m.lowpriority"
	TagServerNotice TagName = "m.server_notice"
)

type Tag struct {
	// Ordering information as a number between 0 and 1.
	// Compared such that 0 is displayed first, and an order of `0.2` would
	// come before a room with order `0.7`.
	//
	// If Order is nil, it should appear last.
	Order *float64 `json:"order"`
}
