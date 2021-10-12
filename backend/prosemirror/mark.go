package prosemirror

// Mark is a Prosemirror document mark.
type Mark struct {
	Type  string                 `json:"type,omitempty"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

// MarkType is a Prosemirror mark specification.
// The ToDOM function is used to render the corresponding mark.
type MarkType struct {
	ToDOM func(mark *Mark, content string) string
}
