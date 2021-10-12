package prosemirror

import (
	"errors"
)

const (
	docType = "doc"
)

// Schema is a Prosemirror schema containing all node and mark specifications.
type Schema struct {
	Nodes map[string]NodeType
	Marks map[string]MarkType
}

// NewSchema returns a new Prosemirror schema for the given nodes and marks.
func NewSchema(nodes map[string]NodeType, marks map[string]MarkType) (*Schema, error) {
	if _, ok := nodes[docType]; !ok {
		return nil, errors.New("schema must have root node")
	}

	return &Schema{nodes, marks}, nil
}
