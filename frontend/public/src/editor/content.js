// Builds and returns table of contents based on headline elements in given DOM.
// It searches recursively for all H2 and H3 text in the tree.
// The second parameter must not be passed, that's used for the recursive call.
export function buildTableOfContents(dom, data) {
	let toc = [];
	let id = 0;
	let h2Number = 0;
	let h3Number = 0;
	let h4Number = 0;
	let activeHeadline = null;

	if(data) {
		id = data.id;
		h2Number = data.h2Number;
		h3Number = data.h3Number;
		h4Number = data.h4Number;
		activeHeadline = data.activeHeadline;
	}

	dom.childNodes.forEach(node => {
		if(!node.tagName) {
			return;
		}

		let tag = node.tagName.toUpperCase();
		let text = node.innerText || "-";

		if(tag === "H2") {
			id++;
			h2Number++;
			toc.push({id, e: node, tag, text: text, number: h2Number, active: false, children: []});
			activeHeadline = toc[toc.length-1];
		}
		else if(tag === "H3") {
			id++;
			h3Number++;

			if(activeHeadline) {
				activeHeadline.children.push({id, e: node, tag, text: text, number: h3Number, active: false});
			}
			else {
				toc.push({id, e: node, tag, text: text, number: h3Number, active: false});
			}
		}
		else if(tag === "H4") {
			id++;
			h4Number++;

			if(activeHeadline) {
				activeHeadline.children.push({id, e: node, tag, text: text, number: h4Number, active: false});
			}
			else {
				toc.push({id, e: node, tag, text: text, number: h4Number, active: false});
			}
		}
		else if(node.childNodes && node.childNodes.length) {
			let result = buildTableOfContents(node, {id, h2Number, h3Number, h4Number, activeHeadline});
			toc = toc.concat(result.toc);
			id = result.data.id;
			h2Number = result.data.h2Number;
			h3Number = result.data.h3Number;
			h4Number = result.data.h4Number;
			activeHeadline = result.data.activeHeadline;
		}
	});

	if(data) {
		return {toc, data: {id, h2Number, h3Number, h4Number, activeHeadline}};
	}

	return toc;
}
