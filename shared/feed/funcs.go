package feed

import (
	"github.com/emvi/hide"
	"github.com/gosimple/slug"
	"text/template"
)

var (
	renderFuncs = template.FuncMap{
		"IdToString": hide.ToString,
		"Add":        Add,
		"SlugWithId": SlugWithId,
	}
)

func Add(a, b int) int {
	return a + b
}

func SlugWithId(str string, id hide.ID) string {
	idStr, _ := hide.ToString(id)
	return slug.Make(str) + "-" + idStr
}
