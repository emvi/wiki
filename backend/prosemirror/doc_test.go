package prosemirror

import (
	"emviwiki/shared/testutil"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestParseDocSimple(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatalf("Document must be parsed, but was: %v", err)
	}

	out, err := json.Marshal(doc)

	if err != nil {
		t.Fatalf("Document must be marshalled, but was: %v", err)
	}

	t.Log(string(out))
	testutil.AssertJSONEquals(t, string(out), sampleSimpleDoc)
}

func TestParseDocComplex(t *testing.T) {
	doc, err := ParseDoc(sampleComplexDoc)

	if err != nil {
		t.Fatalf("Document must be parsed, but was: %v", err)
	}

	out, err := json.Marshal(doc)

	if err != nil {
		t.Fatalf("Document must be marshalled, but was: %v", err)
	}

	t.Log(string(out))
	testutil.AssertJSONEquals(t, string(out), sampleComplexDoc)
}

func TestRenderDocSimple(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatal(err)
	}

	schema, err := NewSchema(map[string]NodeType{
		"doc": {
			ToDOM: func(node *Node, content string) string {
				return content
			},
		},
		"text": {
			ToDOM: func(node *Node, content string) string {
				return node.Text
			},
		},
		"paragraph": {
			ToDOM: func(node *Node, content string) string {
				return fmt.Sprintf("<p>%s</p>", content)
			},
		},
	}, map[string]MarkType{
		"italic": {
			ToDOM: func(mark *Mark, content string) string {
				return fmt.Sprintf("<em>%s</em>", content)
			},
		},
	})

	if err != nil {
		t.Fatalf("Schema must be created, but was: %v", err)
	}

	out, err := RenderDoc(schema, doc)

	if err != nil {
		t.Fatalf("Document must be rendered, but was: %v", err)
	}

	t.Log(out)

	if out != "<p>Hallo <em>Welt</em>!</p>" {
		t.Fatalf("Unexpected HTML: %v", out)
	}
}

func TestRenderDocAttributes(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDocWithAttributes)

	if err != nil {
		t.Fatal(err)
	}

	schema, err := NewSchema(map[string]NodeType{
		"doc": {
			ToDOM: func(node *Node, content string) string {
				return content
			},
		},
		"text": {
			ToDOM: func(node *Node, content string) string {
				return node.Text
			},
		},
		"paragraph": {
			ToDOM: func(node *Node, content string) string {
				return fmt.Sprintf("<p>%s</p>", content)
			},
		},
	}, map[string]MarkType{
		"link": {
			ToDOM: func(mark *Mark, content string) string {
				href := mark.Attrs["href"]
				return fmt.Sprintf(`<a href="%v">%s</a>`, href, content)
			},
		},
	})

	if err != nil {
		t.Fatalf("Schema must be created, but was: %v", err)
	}

	out, err := RenderDoc(schema, doc)

	if err != nil {
		t.Fatalf("Document must be rendered, but was: %v", err)
	}

	t.Log(out)

	if out != `<p><a href="https://emvi.com/">Test</a></p>` {
		t.Fatalf("Unexpected HTML: %v", out)
	}
}

func TestRenderDocUnknownNode(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatal(err)
	}

	schema, err := NewSchema(map[string]NodeType{
		"doc": {
			ToDOM: func(node *Node, content string) string {
				return content
			},
		},
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	out, err := RenderDoc(schema, doc)

	if out != "" || err == nil || err.Error() != "unknown node type 'text'" {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRenderDocUnknownMark(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatal(err)
	}

	schema, err := NewSchema(map[string]NodeType{
		"doc": {
			ToDOM: func(node *Node, content string) string {
				return content
			},
		},
		"text": {
			ToDOM: func(node *Node, content string) string {
				return node.Text
			},
		},
		"paragraph": {
			ToDOM: func(node *Node, content string) string {
				return fmt.Sprintf("<p>%s</p>", content)
			},
		},
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	out, err := RenderDoc(schema, doc)

	if out != "" || err == nil || err.Error() != "unknown mark type 'italic'" {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestFindNodes(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatal(err)
	}

	results := FindNodes(doc, -1, "text", "paragraph")

	if len(results) != 4 {
		t.Fatalf("Expected 4 nodes to be found, but was: %v", len(results))
	}

	for _, result := range results {
		if result.Type != "text" && result.Type != "paragraph" {
			t.Fatalf("Result must be of type 'text' or 'paragraph', but was: %v", result.Type)
		}
	}
}

func TestFindNodesWithLimit(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatal(err)
	}

	results := FindNodes(doc, 2, "text", "paragraph")

	if len(results) != 2 {
		t.Fatalf("Expected 2 nodes to be found, but was: %v", len(results))
	}

	for _, result := range results {
		if result.Type != "text" && result.Type != "paragraph" {
			t.Fatalf("Result must be of type 'text' or 'paragraph', but was: %v", result.Type)
		}
	}
}

func TestTransformNodes(t *testing.T) {
	doc, err := ParseDoc(sampleSimpleDoc)

	if err != nil {
		t.Fatal(err)
	}

	TransformNodes(doc, "text", func(node *Node) {
		node.Text = strings.ToUpper(node.Text)
	})

	results := FindNodes(doc, -1, "text")

	if len(results) != 3 {
		t.Fatalf("Expected 3 nodes to be found, but was: %v", len(results))
	}

	if results[0].Text != "HALLO " ||
		results[1].Text != "WELT" ||
		results[2].Text != "!" {
		t.Fatalf("Expected all texts to be uppercase, but was: %v %v %v", results[0].Text, results[1].Text, results[2].Text)
	}
}

func TestEscapeText(t *testing.T) {
	doc, err := ParseDoc(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hallo "},{"type":"text","marks":[{"type":"italic"}],"text":"<b>Welt</b>"},{"type":"text","text":"!"}]}]}`)

	if err != nil {
		t.Fatal(err)
	}

	EscapeText(doc)
	textNodes := FindNodes(doc, -1, "text")

	if len(textNodes) != 3 {
		t.Fatalf("Three text nodes must be found, but was: %v", len(textNodes))
	}

	if textNodes[0].Text != "Hallo " ||
		textNodes[1].Text != "&lt;b&gt;Welt&lt;/b&gt;" ||
		textNodes[2].Text != "!" {
		t.Fatalf("Text node content not as expected: %v %v %v", textNodes[0].Text, textNodes[1].Text, textNodes[2].Text)
	}
}

const (
	sampleSimpleDoc               = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hallo "},{"type":"text","marks":[{"type":"italic"}],"text":"Welt"},{"type":"text","text":"!"}]}]}`
	sampleSimpleDocWithAttributes = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","marks":[{"type":"link","attrs":{"href":"https://emvi.com/"}}],"text":"Test"}]}]}`
	sampleComplexDoc              = `
{
   "type":"doc",
   "content":[
      {
         "type":"headline",
         "attrs":{
            "level":2
         },
         "content":[
            {
               "type":"text",
               "text":"a"
            }
         ]
      },
      {
         "type":"headline",
         "attrs":{
            "level":3
         },
         "content":[
            {
               "type":"text",
               "text":"b"
            }
         ]
      },
      {
         "type":"paragraph",
         "content":[
            {
               "type":"text",
               "text":"c"
            },
            {
               "type":"text",
               "marks":[
                  {
                     "type":"bold"
                  }
               ],
               "text":"d"
            },
            {
               "type":"text",
               "marks":[
                  {
                     "type":"italic"
                  }
               ],
               "text":"e"
            },
            {
               "type":"text",
               "marks":[
                  {
                     "type":"strikethrough"
                  }
               ],
               "text":"f"
            },
            {
               "type":"text",
               "marks":[
                  {
                     "type":"underlined"
                  }
               ],
               "text":"g"
            }
         ]
      },
      {
         "type":"paragraph",
         "content":[
            {
               "type":"text",
               "marks":[
                  {
                     "type":"link",
                     "attrs":{
                        "href":"https://emvi.com/"
                     }
                  }
               ],
               "text":"h"
            }
         ]
      },
      {
         "type":"bullet_list",
         "content":[
            {
               "type":"list_item",
               "content":[
                  {
                     "type":"paragraph",
                     "content":[
                        {
                           "type":"text",
                           "text":"i"
                        }
                     ]
                  }
               ]
            },
            {
               "type":"list_item",
               "content":[
                  {
                     "type":"paragraph",
                     "content":[
                        {
                           "type":"text",
                           "text":"j"
                        }
                     ]
                  },
                  {
                     "type":"bullet_list",
                     "content":[
                        {
                           "type":"list_item",
                           "content":[
                              {
                                 "type":"paragraph",
                                 "content":[
                                    {
                                       "type":"text",
                                       "text":"k"
                                    }
                                 ]
                              }
                           ]
                        },
                        {
                           "type":"list_item",
                           "content":[
                              {
                                 "type":"paragraph",
                                 "content":[
                                    {
                                       "type":"text",
                                       "text":"l"
                                    }
                                 ]
                              }
                           ]
                        },
                        {
                           "type":"list_item",
                           "content":[
                              {
                                 "type":"paragraph",
                                 "content":[
                                    {
                                       "type":"text",
                                       "text":"m"
                                    }
                                 ]
                              }
                           ]
                        }
                     ]
                  }
               ]
            }
         ]
      },
      {
         "type":"horizontal_rule"
      },
      {
         "type":"image",
         "attrs":{
            "src":"http://localhost:4003/api/v1/content/DoB9mwd3ZV.png"
         }
      },
      {
         "type":"ordered_list",
         "attrs":{
            "order":1
         },
         "content":[
            {
               "type":"list_item",
               "content":[
                  {
                     "type":"paragraph",
                     "content":[
                        {
                           "type":"text",
                           "text":"n"
                        }
                     ]
                  }
               ]
            },
            {
               "type":"list_item",
               "content":[
                  {
                     "type":"paragraph",
                     "content":[
                        {
                           "type":"text",
                           "text":"o"
                        }
                     ]
                  }
               ]
            },
            {
               "type":"list_item",
               "content":[
                  {
                     "type":"paragraph",
                     "content":[
                        {
                           "type":"text",
                           "text":"p"
                        }
                     ]
                  }
               ]
            }
         ]
      },
      {
         "type":"blockquote",
         "content":[
            {
               "type":"paragraph",
               "content":[
                  {
                     "type":"text",
                     "text":"q"
                  }
               ]
            },
            {
               "type":"blockquote",
               "content":[
                  {
                     "type":"paragraph",
                     "content":[
                        {
                           "type":"text",
                           "text":"r"
                        }
                     ]
                  }
               ]
            }
         ]
      },
      {
         "type":"paragraph",
         "content":[
            {
               "type":"file",
               "attrs":{
                  "file":"http://localhost:4003/api/v1/content/FfbLp3jwy0",
                  "name":"12QnJATs51WK_49mb",
                  "size":"49.28 MB"
               }
            }
         ]
      }
   ]
}
`
)
