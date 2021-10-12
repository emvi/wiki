import {Plugin} from "prosemirror-state";

export class EditorMenu {
	init(view, items) {
		this.view = view;
		this.buildDOM(items);

		this.mouseupEventHandler = e => {this.mouseUp(e);};
		window.addEventListener("mouseup", this.mouseupEventHandler);
	}

	destroy() {}

	buildDOM(items) {
		this.dom = document.createElement("div");
		this.dom.className = "bottom-menu";
		this.dom.style.display = "none";
		this.items = [];

		for(let item of items) {
			item.dom = EditorMenu.buildItem(item.action.tooltip, item.action.icon);
			this.dom.appendChild(item.dom);
			item.dom.addEventListener("mousedown", e => {EditorMenu.mouseDown(e, this.view, item);});
			this.items.push(item);
		}
	}

	static buildItem(tooltip, icon) {
		let i = document.createElement("i");
		i.className = "icon "+icon;
		i.title = tooltip;
		return i;
	}

	update() {
		if(this.view.dom.classList.contains("ProseMirror-focused")) {
			this.enableItems();
			this.dom.style.display = "";
		}
		else {
			this.dom.style.display = "none";
		}
	}

	enableItems() {
		this.items.forEach(({command, dom}) => {
			let active = command(this.view.state, null, this.view);

			if(active) {
				dom.classList.remove("disabled");
			}
			else {
				dom.classList.add("disabled");
			}
		});
	}

	mouseUp() {
		this.update();
	}

	static mouseDown(e, view, item) {
		e.preventDefault();
		view.focus();
		item.command(view.state, view.dispatch, view);
	}

	getDOM() {
		return this.dom;
	}
}

export function menuPlugin(items) {
	let menu = new EditorMenu();
	let plugin = new Plugin({
		view(view) {
			if(!this.menu) {
				menu.init(view, items);
				this.menu = menu;
				view.dom.parentNode.insertBefore(this.menu.getDOM(), view.dom);
			}

			return this.menu;
		}
	});

	plugin.destroy = function() {
		menu.destroy();
	};

	return plugin;
}
