package schema

import (
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/feed"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
)

func GetMarkdownSchema(orga *model.Organization) *prosemirror.Schema {
	nodes := map[string]prosemirror.NodeType{
		"text": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return node.Text
			},
		},
		"hard_break": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return "\n"
			},
		},
		"paragraph": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				if content == "" {
					return "\n"
				}

				return fmt.Sprintf("%s\n", content)
			},
		},
		"blockquote": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				lines := strings.Split(content, "\n")
				out := ""

				for i, line := range lines {
					if i < len(lines)-1 {
						out += fmt.Sprintf("> %s\n", line)
					}
				}

				return fmt.Sprintf("\n%s\n", out)
			},
		},
		"code_block": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				language := node.Attrs["language"]

				if language != nil {
					return fmt.Sprintf("```%s\n%s\n```", language, content)
				}

				return fmt.Sprintf("```\n%s\n```", content)
			},
		},
		"infobox": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				lines := strings.Split(content, "\n")
				out := ""

				for i, line := range lines {
					if i < len(lines)-1 {
						out += fmt.Sprintf("> %s\n", line)
					}
				}

				return fmt.Sprintf("\n%s\n", out)
			},
		},
		"headline": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				level, ok := node.Attrs["level"].(float64)

				if ok {
					if level == 2 {
						return fmt.Sprintf("\n## %s\n\n", content)
					} else if level == 3 {
						return fmt.Sprintf("\n### %s\n\n", content)
					} else if level == 4 {
						return fmt.Sprintf("\n#### %s\n\n", content)
					}
				}

				return fmt.Sprintf("%s\n\n", content)
			},
		},
		"horizontal_rule": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return "\n---\n\n"
			},
		},
		"image": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]

				if content == "" {
					return fmt.Sprintf("![%v](%v)\n", src, src)
				}

				return fmt.Sprintf("![%v](%v) (%s)\n", src, src, content[:len(content)-1])
			},
		},
		"file": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				file := node.Attrs["file"]
				name := node.Attrs["name"]
				size := node.Attrs["size"]
				return fmt.Sprintf(`[%v](%s) (%v)`, name, file, size)
			},
		},
		"mention": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				id := node.Attrs["id"].(string)
				t := node.Attrs["type"].(string)
				title := node.Attrs["title"].(string)

				objId, _ := hide.FromString(id)
				links := map[string]string{
					"article": "/read/" + feed.SlugWithId(title, objId),
					"list":    "/list/" + feed.SlugWithId(title, objId),
					"user":    "/member/" + id,
					"tag":     "/tag/" + id,
					"group":   "/group/" + feed.SlugWithId(title, objId),
				}
				link := util.GetOrganizationURL(orga, links[t])

				return fmt.Sprintf(`@[%s](%s)`, title, link)
			},
		},
		"ordered_list": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				lines := strings.Split(content, "\n")
				lines = lines[:len(lines)-1]
				index := 1
				out := ""

				for _, line := range lines {
					if strings.HasPrefix(line, "- ") {
						out += fmt.Sprintf("%d.%s\n", index, line[1:])
						index++
					} else {
						out += fmt.Sprintf(" %s\n", line)
					}
				}

				return out
			},
		},
		"bullet_list": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return content
			},
		},
		"list_item": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				lines := strings.Split(content, "\n")
				lines = lines[:len(lines)-1]
				out := ""

				for i, line := range lines {
					if i > 0 {
						out += fmt.Sprintf("  %s\n", line)
					} else {
						out += fmt.Sprintf("%s\n", line)
					}
				}

				return fmt.Sprintf("- %s", out)
			},
		},
		"check_list": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return content
			},
		},
		"check_list_item": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				lines := strings.Split(content, "\n")
				lines = lines[:len(lines)-1]
				out := ""

				for i, line := range lines {
					if i > 0 {
						out += fmt.Sprintf("  %s\n", line)
					} else {
						out += fmt.Sprintf("%s\n", line)
					}
				}

				checked := node.Attrs["checked"].(bool)

				if checked {
					return fmt.Sprintf("- [x] %s", out)
				}

				return fmt.Sprintf("- [ ] %s", out)
			},
		},
		"table": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				lines := strings.Split(content, "\n")
				out := ""

				for i, line := range lines {
					out += fmt.Sprintf("%s\n", line)

					if i == 0 {
						columns := strings.Count(line, "|") - 1

						for i := 0; i < columns; i++ {
							out += "| --- "
						}

						out += "|\n"
					}
				}

				return fmt.Sprintf("\n%s", out)
			},
		},
		"table_row": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf("%s|\n", content)
			},
		},
		"table_cell": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				if strings.HasSuffix(content, "\n") {
					content = content[:len(content)-1]
				}

				return fmt.Sprintf(`| %s `, content)
			},
		},
		"table_header": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				if strings.HasSuffix(content, "\n") {
					content = content[:len(content)-1]
				}

				return fmt.Sprintf(`| %s `, content)
			},
		},
		"youtube": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf("[%s](%s)\n", src, src)
			},
		},
		"vimeo": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf("[%s](%s)\n", src, src)
			},
		},
		"spotify": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf("[%s](%s)\n", src, src)
			},
		},
		"pdf": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf("[%s](%s)\n", src, src)
			},
		},
		"link_preview": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				href := node.Attrs["href"]
				title := node.Attrs["title"]
				return fmt.Sprintf("[%s](%s)\n", title, href)
			},
		},
		"doc": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return content
			},
		},
	}
	marks := map[string]prosemirror.MarkType{
		"bold": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("**%s**", content)
			},
		},
		"italic": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("*%s*", content)
			},
		},
		"underlined": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("_%s_", content)
			},
		},
		"strikethrough": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("~~%s~~", content)
			},
		},
		"code": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("`%s`", content)
			},
		},
		"link": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				href := mark.Attrs["href"]
				return fmt.Sprintf(`[%v](%v)`, content, href)
			},
		},
		"sub": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("<sub>%s</sub>", content)
			},
		},
		"sup": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("<sup>%s</sup>", content)
			},
		},
	}
	s, err := prosemirror.NewSchema(nodes, marks)

	if err != nil {
		logbuch.Fatal("Error build article schema", logbuch.Fields{"err": err})
	}

	return s
}
