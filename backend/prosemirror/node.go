package prosemirror

// Node is a Prosemirror document node.
type Node struct {
	Type    string                 `json:"type,omitempty"`
	Attrs   map[string]interface{} `json:"attrs,omitempty"`
	Content []Node                 `json:"content,omitempty"`
	Marks   []Mark                 `json:"marks,omitempty"`
	Text    string                 `json:"text,omitempty"`
}

// NodeType is a Prosemirror node specification.
// The ToDOM function is used to render the corresponding node.
type NodeType struct {
	ToDOM func(node *Node, content string) string
}
