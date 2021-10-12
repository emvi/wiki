const {Schema} = require("prosemirror-model");
const {addListNodes} = require("prosemirror-schema-list");
const {tableNodes} = require("prosemirror-tables");

const nodes = {
	dummy_table: {
		group: "table" // register group
	},
	text: {
		group: "inline"
	},
	hard_break: {
		group: "inline",
		inline: true,
		selectable: false
	},
	paragraph: {
		group: "block",
		content: "inline*"
	},
	blockquote: {
		group: "block",
		content: "(block|table)+"
	},
	code_block: {
		group: "block",
		content: "text*",
		marks: "",
		attrs: {language: {default: "text/plain"}},
		code: true,
		defining: true
	},
	infobox: {
		group: "block",
		content: "(block|table)+",
		draggable: true,
		attrs: {
			color: {default: "blue"}
		}
	},
	headline: {
		group: "block",
		content: "inline*",
		attrs: {level: {default: 2}},
		defining: true
	},
	horizontal_rule: {
		group: "block",
		marks: ""
	},
	image: {
		group: "block",
		content: "paragraph",
		atom: true,
		draggable: true,
		attrs: {src: {}}
	},
	file: {
		group: "inline",
		inline: true,
		draggable: true,
		attrs: {file: {}, name: {}, size: {}}
	},
	mention: {
		group: "inline",
		inline: true,
		draggable: true,
		attrs: {id: {}, type: {}, title: {}, time: {}}
	},
	check_list: {
		group: "block",
		content: "check_list_item+"
	},
	check_list_item: {
		content: "paragraph block*",
		attrs: {checked: {default: false}}
	},
	youtube: {
		group: "block",
		marks: "",
		attrs: {src: {}}
	},
	vimeo: {
		group: "block",
		marks: "",
		attrs: {src: {}}
	},
	spotify: {
		group: "block",
		marks: "",
		attrs: {src: {}}
	},
	pdf: {
		group: "block",
		marks: "",
		attrs: {src: {}}
	},
	link_preview: {
		group: "block",
		marks: "",
		attrs: {href: {}, title: {default: ""}, description: {default: ""}, image: {default: ""}}
	},
	doc: {
		content: "(block|table)+"
	}
};

const marks = {
	bold: {},
	italic: {},
	underlined: {},
	strikethrough: {},
	code: {},
	link: {
		attrs: {href: {}},
		inclusive: false
	},
	sub: {
		toDOM() {return ["sub", 0]},
		parseDOM: [{tag: "sub"}]
	},
	sup: {
		toDOM() {return ["sup", 0]},
		parseDOM: [{tag: "sup"}]
	}
};

function addTableNodes(nodes) {
	return nodes.append(tableNodes({
		tableGroup: "table",
		cellContent: "block+",
		cellAttributes: {
			background: {
				default: "none",
				getFromDOM: dom => {
					return dom.style.background || null;
				},
				setDOMAttr: (value, attrs) => {
					if(value && value !== "none") {
						attrs.style = (attrs.style || "") + `background: ${value};`
					}
				}
			}
		}
	}));
}

let schema = new Schema({nodes, marks});
module.exports.schema = new Schema({
	nodes: addTableNodes(addListNodes(schema.spec.nodes, "paragraph block*", "block")),
	marks: schema.spec.marks
});
