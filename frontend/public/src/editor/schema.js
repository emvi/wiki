import {Schema} from "prosemirror-model";
import {addListNodes} from "prosemirror-schema-list";
import {tableNodes} from "prosemirror-tables";

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
		selectable: false,
		toDOM() {return ["br"]},
		parseDOM: [{tag: "br"}]
	},
	paragraph: {
		group: "block",
		content: "inline*",
		toDOM() {return ["p", 0]},
		parseDOM: [{tag: "p"}]
	},
	blockquote: {
		group: "block",
		content: "(block|table)+",
		toDOM() {return ["blockquote", 0]},
		parseDOM: [{tag: "blockquote"}]
	},
	code_block: {
		group: "block",
		content: "text*",
		marks: "",
		attrs: {language: {default: "text/plain"}},
		code: true,
		defining: true,
		toDOM(node) {return ["pre", ["code", node.attrs, 0]]},
		parseDOM: [{
			tag: "pre",
			preserveWhitespace: "full",
			getAttrs(dom) {
				return {language: dom.getAttribute("language")};
			}
		}]
	},
	infobox: {
		group: "block",
		content: "(block|table)+",
		draggable: true,
		attrs: {
			color: {default: "blue"}
		},
		toDOM(node) {
			return [
				"div",
				{class: `infobox ${node.attrs.color}`, color: node.attrs.color},
				0
			];
		},
		parseDOM: [{
			tag: "div.infobox",
			getAttrs(dom) {
				return {color: dom.getAttribute("color")};
			}
		}]
	},
	headline: {
		group: "block",
		content: "inline*",
		attrs: {level: {default: 2}},
		defining: true,
		toDOM(node){return ["h" + node.attrs.level, 0]},
		parseDOM: [
			{tag: "h2", attrs: {level: 2}},
			{tag: "h3", attrs: {level: 3}},
			{tag: "h4", attrs: {level: 4}}
		]
	},
	horizontal_rule: {
		group: "block",
		marks: "",
		toDOM() {return ["hr"]},
		parseDOM: [{tag: "hr"}]
	},
	image: {
		group: "block",
		content: "paragraph",
		atom: true,
		draggable: true,
		attrs: {src: {}},
		toDOM(node) {
			return [
				"figure",
				{src: node.attrs.src},
				["img", {src: node.attrs.src, alt: node.attrs.src}],
				["figcaption", 0]
			]
		},
		parseDOM: [{
			tag: "figure",
			getAttrs(dom) {
				return {src: dom.getAttribute("src")};
			}
		}]
	},
	file: {
		group: "inline",
		inline: true,
		draggable: true,
		attrs: {file: {}, name: {}, size: {}},
		toDOM(node) {
			let dom = document.createElement("template");
			dom.innerHTML = `<a href="${node.attrs.file}" class="file" file="${node.attrs.file}" name="${node.attrs.name}" size="${node.attrs.size}" title="${node.attrs.name}"><i class="icon icon-file"></i> <span class="name">${node.attrs.name}</span> <span class="size">(${node.attrs.size})</span></a>`;
			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "a[file]",
			getAttrs(dom) {
				return {
					file: dom.getAttribute("file"),
					name: dom.getAttribute("name"),
					size: dom.getAttribute("size")
				};
			}
		}]
	},
	mention: {
		group: "inline",
		inline: true,
		draggable: true,
		attrs: {id: {}, type: {}, title: {}, time: {}},
		toDOM(node) {
			let dom = document.createElement("template");
			dom.innerHTML = `<span class="${node.attrs.type}" mention="${node.attrs.id}" object="${node.attrs.type}" title="${node.attrs.title}" time="${node.attrs.time}">${node.attrs.title}</span>`;
			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "span[mention]",
			getAttrs(dom) {
				return {
					id: dom.getAttribute("mention"),
					type: dom.getAttribute("object"),
					title: dom.getAttribute("title"),
					time: dom.getAttribute("time")
				};
			}
		}]
	},
	check_list: {
		group: "block",
		content: "check_list_item+",
		toDOM() {
			return [
				"ul",
				{checklist: "", class: "checklist"},
				0
			];
		},
		parseDOM: [{
			tag: "ul[checklist]"
		}]
	},
	check_list_item: {
		content: "paragraph block*",
		attrs: {checked: {default: false}},
		toDOM(node) {
			return [
				"li",
				{checklistitem: "", class: node.attrs.checked ? "checked" : "", checked: node.attrs.checked.toString()},
				["span", 0],
				["span", {contenteditable: false, class: "checklist-trigger"}]
			];
		},
		parseDOM: [{
			tag: "li[checklistitem]",
			getAttrs(dom) {
				return {checked: dom.getAttribute("checked") === "true"};
			}
		}]
	},
	youtube: {
		group: "block",
		marks: "",
		attrs: {src: {}},
		toDOM(node) {
			let dom = document.createElement("template");
			dom.innerHTML = `<div class="embed youtube" data-src="${node.attrs.src}"><iframe src="${node.attrs.src}" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen youtube></iframe></div>`;
			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "div.youtube",
			getAttrs(dom) {
				return {src: dom.getAttribute("data-src")}
			}
		}]
	},
	vimeo: {
		group: "block",
		marks: "",
		attrs: {src: {}},
		toDOM(node) {
			let dom = document.createElement("template");
			dom.innerHTML = `<div class="embed vimeo" data-src="${node.attrs.src}"><iframe title="vimeo-player" src="${node.attrs.src}" frameborder="0" allowfullscreen vimeo></iframe></div>`;
			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "div.vimeo",
			getAttrs(dom) {
				return {src: dom.getAttribute("data-src")}
			}
		}]
	},
	spotify: {
		group: "block",
		marks: "",
		attrs: {src: {}},
		toDOM(node) {
			let dom = document.createElement("template");
			dom.innerHTML = `<div class="embed spotify" data-src="${node.attrs.src}"><iframe src="${node.attrs.src}" frameborder="0" allowtransparency="true" allow="encrypted-media" spotify></iframe></div>`;
			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "div.spotify",
			getAttrs(dom) {
				return {src: dom.getAttribute("data-src")}
			}
		}]
	},
	pdf: {
		group: "block",
		marks: "",
		attrs: {src: {}},
		toDOM(node) {
			let dom = document.createElement("template");
			dom.innerHTML = `<div class="embed pdf" data-src="${node.attrs.src}"><iframe src="${node.attrs.src}" frameborder="0" allowtransparency="true" allow="encrypted-media" pdf></iframe></div>`;
			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "div.pdf",
			getAttrs(dom) {
				return {src: dom.getAttribute("data-src")}
			}
		}]
	},
	link_preview: {
		group: "block",
		marks: "",
		attrs: {href: {}, title: {default: ""}, description: {default: ""}, image: {default: ""}},
		toDOM(node) {
			let dom = document.createElement("template");

			if(node.attrs.image) {
				dom.innerHTML = `<a href="${node.attrs.href}" target="_blank" class="embed link" data-href="${node.attrs.href}" data-title="${node.attrs.title}" data-description="${node.attrs.description}" data-image="${node.attrs.image}"><div class="image"><img src="${node.attrs.image}" alt="" /></div><div class="info"><div class="title">${node.attrs.title}</div><div class="description">${node.attrs.description}</div><div class="url">${node.attrs.href}</div></div></a>`;
			}
			else {
				dom.innerHTML = `<a href="${node.attrs.href}" target="_blank" class="embed link" data-href="${node.attrs.href}" data-title="${node.attrs.title}" data-description="${node.attrs.description}" data-image=""><div class="info"><div class="title">${node.attrs.title}</div><div class="description">${node.attrs.description}</div><div class="url">${node.attrs.href}</div></div></a>`;
			}

			return dom.content.firstChild;
		},
		parseDOM: [{
			tag: "div.link",
			getAttrs(dom) {
				return {
					href: dom.getAttribute("data-href"),
					title: dom.getAttribute("data-title"),
					description: dom.getAttribute("data-description"),
					image: dom.getAttribute("data-image")
				}
			}
		}]
	},
	doc: {
		content: "(block|table)+"
	}
};

const marks = {
	bold: {
		toDOM() {return ["strong", 0]},
		parseDOM: [{tag: "strong"}]
	},
	italic: {
		toDOM() {return ["em", 0]},
		parseDOM: [{tag: "em"}]
	},
	underlined: { // highlight
		toDOM() {return ["u", 0]},
		parseDOM: [{tag: "u"}]
	},
	strikethrough: {
		toDOM() {return ["strike", 0]},
		parseDOM: [{tag: "strike"}]
	},
	code: {
		toDOM() {return ["code", 0]},
		parseDOM: [{tag: "code"}]
	},
	link: {
		attrs: {href: {}},
		inclusive: false,
		toDOM(node) {
			return ["a", {href: node.attrs.href, target: "_blank", rel: "noreferrer"}, 0];
		},
		parseDOM: [{
			tag: "a[href]",
			getAttrs(dom) {
				return {href: dom.getAttribute("href")};
			}
		}]
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
export let editorSchema = new Schema({
	nodes: addTableNodes(addListNodes(schema.spec.nodes, "paragraph block*", "block")),
	marks: schema.spec.marks
});
