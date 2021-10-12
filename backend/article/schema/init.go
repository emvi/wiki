package schema

import (
	"emviwiki/backend/prosemirror"
)

var (
	HTMLSchema *prosemirror.Schema
)

func init() {
	buildHTMLSchema()
}
