package schema

import (
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/feed"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strconv"
)

func buildHTMLSchema() {
	nodes := map[string]prosemirror.NodeType{
		"text": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return node.Text
			},
		},
		"hard_break": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return "<br />"
			},
		},
		"paragraph": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				if content == "" {
					content = "<br />"
				}

				return fmt.Sprintf("<p>%s</p>", content)
			},
		},
		"blockquote": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf("<blockquote>%s</blockquote>", content)
			},
		},
		"code_block": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				language := node.Attrs["language"]

				if language != nil {
					return fmt.Sprintf(`<pre><code language="%s">%s</code></pre>`, language, content)
				}

				return fmt.Sprintf(`<pre><code>%s</code></pre>`, content)
			},
		},
		"infobox": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				color := node.Attrs["color"]
				return fmt.Sprintf(`<div class="infobox %s" color="%s">%s</div>`, color, color, content)
			},
		},
		"headline": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				level := node.Attrs["level"]
				return fmt.Sprintf("<h%v>%s</h%v>", level, content, level)
			},
		},
		"horizontal_rule": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return "<hr />"
			},
		},
		"image": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]

				if len(node.Content) == 0 ||
					(len(node.Content) == 1 &&
						(len(node.Content[0].Content) == 0 || (len(node.Content[0].Content) == 1 && node.Content[0].Content[0].Text == ""))) {
					return fmt.Sprintf(`<img src="%v" alt="%v" />`, src, src)
				}

				return fmt.Sprintf(`<figure><img src="%v" alt="%v" /><figcaption>%s</figcaption></figure>`, src, src, content)
			},
		},
		"file": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				file := node.Attrs["file"]
				name := node.Attrs["name"]
				size := node.Attrs["size"]
				return fmt.Sprintf(`<a href="%s" class="file" file="%v" name="%v" size="%v" title="%v"><i class="icon icon-file"></i> <span class="name">%v</span> <span class="size">(%v)</span></></a>`, file, file, name, size, name, name, size)
			},
		},
		"mention": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				id := node.Attrs["id"].(string)
				t := node.Attrs["type"].(string)
				title := node.Attrs["title"].(string)
				time := node.Attrs["time"]

				objId, _ := hide.FromString(id)
				links := map[string]string{
					"article": "/read/" + feed.SlugWithId(title, objId),
					"list":    "/list/" + feed.SlugWithId(title, objId),
					"user":    "/member/" + id,
					"tag":     "/tag/" + id,
					"group":   "/group/" + feed.SlugWithId(title, objId),
				}
				link := links[t]

				// render as link so it can be used for clients, this differs from the collab/frontend implementation!
				return fmt.Sprintf(`<a href="%s" class="%s" mention="%s" object="%v" title="%s" time="%v">%s</a>`, link, t, id, t, title, time, title)
			},
		},
		"ordered_list": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf(`<ol>%s</ol>`, content)
			},
		},
		"bullet_list": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf("<ul>%s</ul>", content)
			},
		},
		"list_item": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf("<li>%s</li>", content)
			},
		},
		"check_list": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf(`<ul class="checklist">%s</ul>`, content)
			},
		},
		"check_list_item": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				checked, ok := node.Attrs["checked"].(bool)
				checkboxClass := ""

				if ok && checked {
					checkboxClass = "checked"
				}

				return fmt.Sprintf(`<li class="%s">%s</li>`, checkboxClass, content)
			},
		},
		"table": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf(`<div class="table-wrapper"><div class="table-content"><table><tbody>%s</tbody></table></div></div>`, content)
			},
		},
		"table_row": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				return fmt.Sprintf("<tr>%s</tr>", content)
			},
		},
		"table_cell": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				colspan := node.Attrs["colspan"].(float64)
				rowspan := node.Attrs["rowspan"].(float64)
				colwidth, ok := node.Attrs["colwidth"].([]interface{})
				width := "auto"

				if ok && len(colwidth) >= 1 {
					width = strconv.Itoa(int(colwidth[0].(float64))) + "px"
				}

				background, ok := node.Attrs["background"].(string)

				if !ok || background == "none" {
					background = ""
				} else if ok && background != "none" {
					background = fmt.Sprintf(`background:%s;`, background)
				}

				return fmt.Sprintf(`<td colspan="%v" rowspan="%v" style="width:%s;%s">%s</td>`, int(colspan), int(rowspan), width, background, content)
			},
		},
		"table_header": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				colspan := node.Attrs["colspan"].(float64)
				rowspan := node.Attrs["rowspan"].(float64)
				colwidth, ok := node.Attrs["colwidth"].([]interface{})
				width := "auto"

				if ok && len(colwidth) >= 1 {
					width = strconv.Itoa(int(colwidth[0].(float64))) + "px"
				}

				background, ok := node.Attrs["background"].(string)

				if !ok || background == "none" {
					background = ""
				} else if ok && background != "none" {
					background = fmt.Sprintf(`background:%s;`, background)
				}

				return fmt.Sprintf(`<th colspan="%v" rowspan="%v" style="width:%s;%s">%s</th>`, int(colspan), int(rowspan), width, background, content)
			},
		},
		"youtube": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf(`<div class="embed youtube" data-src="%v"><iframe src="%v" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen youtube></iframe></div>`, src, src)
			},
		},
		"vimeo": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf(`<div class="embed vimeo" data-src="%v"><iframe title="vimeo-player" src="%v" frameborder="0" allowfullscreen vimeo></iframe></div>`, src, src)
			},
		},
		"spotify": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf(`<div class="embed spotify" data-src="%v"><iframe src="%v" frameborder="0" allowtransparency="true" allow="encrypted-media" spotify></iframe></div>`, src, src)
			},
		},
		"pdf": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				src := node.Attrs["src"]
				return fmt.Sprintf(`<div class="embed pdf" data-src="%v"><iframe src="%v" frameborder="0" allowtransparency="true" allow="encrypted-media" pdf></iframe></div>`, src, src)
			},
		},
		"link_preview": {
			ToDOM: func(node *prosemirror.Node, content string) string {
				href := node.Attrs["href"]
				title := node.Attrs["title"]
				description := node.Attrs["description"]
				image, ok := node.Attrs["image"].(string)

				if ok && image != "" {
					return fmt.Sprintf(`<a href="%v" target="_blank" class="embed link" data-href="%v" data-title="%v" data-description="%v" data-image="%v"><div class="image"><img src="%v" alt="" /></div><div class="info"><div class="title">%v</div><div class="description">%v</div><div class="url">%v</div></div></a>`, href, href, title, description, image, image, title, description, href)
				}

				return fmt.Sprintf(`<a href="%v" target="_blank" class="embed link" data-href="%v" data-title="%v" data-description="%v" data-image=""><div class="info"><div class="title">%v</div><div class="description">%v</div><div class="url">%v</div></div></a>`, href, href, title, description, title, description, href)
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
				return fmt.Sprintf("<strong>%s</strong>", content)
			},
		},
		"italic": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("<em>%s</em>", content)
			},
		},
		"underlined": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("<u>%s</u>", content)
			},
		},
		"strikethrough": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("<strike>%s</strike>", content)
			},
		},
		"code": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				return fmt.Sprintf("<code>%s</code>", content)
			},
		},
		"link": {
			ToDOM: func(mark *prosemirror.Mark, content string) string {
				href := mark.Attrs["href"]
				return fmt.Sprintf(`<a href="%v" target="_blank" rel="noreferrer">%v</a>`, href, content)
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

	HTMLSchema = s
}
