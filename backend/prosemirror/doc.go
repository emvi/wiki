package prosemirror

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"strings"
)

// TransformFunc is a transformation function to modify nodes inside a document.
type TransformFunc func(*Node)

// ParseDoc parses a Prosemirror JSON document into a node and mark structure.
func ParseDoc(doc string) (*Node, error) {
	node := new(Node)

	if err := json.Unmarshal([]byte(doc), node); err != nil {
		return nil, err
	}

	return node, nil
}

// RenderDoc renders the given Prosemirror node to its HTML representation up to the given character limit.
// To render the full document, a limit less than zero can be used.
// Before rendering, all text is HTML escaped using EscapeText.
// Returns an error if a node or mark specification cannot be found in the schema.
func RenderDoc(schema *Schema, doc *Node) (string, error) {
	EscapeText(doc)
	return renderNode(schema, doc)
}

func renderNode(schema *Schema, node *Node) (string, error) {
	var out strings.Builder

	// render children
	if len(node.Content) != 0 {
		for _, content := range node.Content {
			o, err := renderNode(schema, &content)

			if err != nil {
				return "", err
			}

			out.WriteString(o)
		}
	}

	// find node in schema
	t, ok := schema.Nodes[node.Type]

	if !ok {
		return "", errors.New(fmt.Sprintf("unknown node type '%v'", node.Type))
	}

	// render node and apply marks
	outStr := t.ToDOM(node, out.String())

	for i := len(node.Marks) - 1; i >= 0; i-- {
		o, err := renderMark(schema, &node.Marks[i], outStr)

		if err != nil {
			return "", err
		}

		outStr = o
	}

	return outStr, nil
}

func renderMark(schema *Schema, mark *Mark, content string) (string, error) {
	t, ok := schema.Marks[mark.Type]

	if !ok {
		return "", errors.New(fmt.Sprintf("unknown mark type '%v'", mark.Type))
	}

	return t.ToDOM(mark, content), nil
}

// FindNodes returns up to the set limit nodes of the given types.
// To find all nodes, a limit less than zero can be used.
func FindNodes(doc *Node, limit int, typeNames ...string) []Node {
	typeNamesMap := make(map[string]bool)

	for _, typeName := range typeNames {
		typeNamesMap[strings.ToLower(typeName)] = true
	}

	nodes, _ := findNodesByType(doc, typeNamesMap, limit)
	return nodes
}

func findNodesByType(node *Node, typeNames map[string]bool, limit int) ([]Node, int) {
	results := make([]Node, 0)

	if _, ok := typeNames[strings.ToLower(node.Type)]; ok {
		results = append(results, *node)
		limit--
	}

	if limit == 0 {
		return results, limit
	}

	for _, child := range node.Content {
		nodes, limit := findNodesByType(&child, typeNames, limit)
		results = append(results, nodes...)

		if limit == 0 {
			return results, limit
		}
	}

	return results, limit
}

// TransformNodes transforms all nodes of given type with given transformation function.
func TransformNodes(doc *Node, typeName string, transform TransformFunc) {
	typeName = strings.ToLower(typeName)
	transformNodesByType(doc, typeName, transform)
}

func transformNodesByType(node *Node, typeName string, transform TransformFunc) {
	if strings.ToLower(node.Type) == typeName {
		transform(node)
	}

	for i := range node.Content {
		transformNodesByType(&node.Content[i], typeName, transform)
	}
}

// EscapeText HTML escapes all text within the given document.
func EscapeText(doc *Node) {
	for i := range doc.Content {
		EscapeText(&doc.Content[i])
	}

	doc.Text = html.EscapeString(doc.Text)
}
