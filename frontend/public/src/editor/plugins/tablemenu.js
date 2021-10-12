import {Plugin} from "prosemirror-state";

const EDITOR_WRAPPER_ID = "editor-wrapper";
const CELL_COLORS = [
    "rgba(25,140,255,.2)",
    "rgba(203,69,229,.2)",
    "rgba(255,76,195,.2)",
    "rgba(255,26,26,.2)",
    "rgba(255,140,25,.2)",
    "rgba(255,216,26,.2)",
    "rgba(20,204,82,.2)",
    "none"
];
const MAX_WINDOW_WIDTH = 512;

class TableCmdMenu {
    constructor(menu, view, items) {
        this.menu = menu;
        this.view = view;
        this.items = items;
        this.buildDOM();
    }

    buildDOM() {
        this.dom = document.createElement("div");
        this.dom.className = "table-menu";

        for(let item of this.items) {
            item.dom = TableCmdMenu.buildItem(item.action.tooltip, item.action.icon);
            this.dom.appendChild(item.dom);
            this.handleClick(this.view, item);
        }
    }

    static buildItem(tooltip, icon) {
        let i = document.createElement("i");
        i.className = "icon "+icon;
        i.title = tooltip;
        return i;
    }

    handleClick(view, item) {
        item.dom.addEventListener("mouseup", e => {
            e.preventDefault();
            e.stopPropagation();
            view.focus();

            if(item.state) {
                this.menu.switchState(item.state, item.command);
            }
            else {
                item.command(view.state, view.dispatch);
            }
        });
    }

    enableItems() {
        let activeItemFound = false;

        this.items.forEach(({command, dom}) => {
            let active = command(this.view.state, null, this.view);
            dom.style.display = active ? "" : "none";

            if(active) {
                activeItemFound = true;
            }
        });

        return activeItemFound;
    }

    activate() {}

    reset() {}

    setCommand(command) {}

    getDOM() {
        return this.dom;
    }
}

class TableCellColorMenu {
    constructor(menu, view) {
        this.menu = menu;
        this.view = view;
        this.buildDOM();
    }

    buildDOM() {
        this.dom = document.createElement("div");
        this.dom.className = "table-menu";

        for(let color of CELL_COLORS) {
            let dom = TableCellColorMenu.buildColor(color);
            this.dom.appendChild(dom);
            this.handleClick(this.view, dom, color);
        }

        let cancel = TableCellColorMenu.buildButton("Cancel", "icon-back");
        cancel.addEventListener("mouseup", e => {
            e.preventDefault();
            e.stopPropagation();
            this.menu.switchState("default");
        });
        this.dom.appendChild(cancel);
    }

    static buildColor(color) {
        let div = document.createElement("div");
        div.className = "cell-color";
        div.style.backgroundColor = color;
        return div;
    }

    static buildButton(tooltip, icon) {
        let i = document.createElement("i");
        i.className = "icon "+icon;
        i.title = tooltip;
        return i;
    }

    handleClick(view, dom, color) {
        dom.addEventListener("mouseup", e => {
            e.preventDefault();
            e.stopPropagation();
            view.focus();
            this.command("background", color)(view.state, view.dispatch);
        });
    }

    cancel() {
        this.menu.switchState("default");
    }

    enableItems() {
        return true;
    }

    activate() {}

    reset() {}

    setCommand(command) {
        this.command = command;
    }

    getDOM() {
        return this.dom;
    }
}

export class TableMenu {
    init(view, items) {
        this.view = view;
        this.buildDOM(items);

        this.mouseupEventHandler = e => {this.mouseUp(e);};
        this.scrollEventHandler = () => {this.calcPosition();};
        window.addEventListener("mouseup", this.mouseupEventHandler);
        window.addEventListener("scroll", this.scrollEventHandler);
    }

    destroy() {
        // FIXME gets called somehow...
        //window.removeEventListener("mouseup", this.mouseupEventHandler);
    }

    buildDOM(items) {
        this.states = {
            default: new TableCmdMenu(this, this.view, items),
            color: new TableCellColorMenu(this, this.view)
        };
        this.switchState("default");
    }

    switchState(name, command) {
        if(this.activeState === this.states[name]) {
            return;
        }

        this.activeState = this.states[name];
        this.activeState.reset();
        this.activeState.setCommand(command);
        this.activeState.activate();
        this.replaceDOM();
        this.view.focus();
        this.update(this.view);
    }

    replaceDOM() {
        if(this.dom) {
            let dom = this.activeState.getDOM();
            this.view.dom.parentNode.replaceChild(dom, this.dom);
            this.dom = dom;
        }
        else {
            this.dom = this.activeState.getDOM();
        }
    }

    update(view) {
        let state = view.state;

        if(!view.hasFocus()) {
            this.close();
            return;
        }

        let tablePos = this.getSelectedTablePos(state.selection);

        if(tablePos === -1 || !this.activeState.enableItems()) {
            this.close();
            return;
        }

        this.tableDOM = view.domAtPos(tablePos).node;
        this.calcPosition();
    }

    getSelectedTablePos(selection) {
        let from = selection.$from;

        for(let i = from.depth; i > 0; i--) {
            let node = from.node(i);

            if(node.type.name === "table") {
                return from.start(i);
            }
        }

        return -1;
    }

    calcPosition() {
        if(!this.tableDOM) {
            return;
        }

        this.dom.style.display = ""; // display before getting width or else it returns 0

        if(window.innerWidth < MAX_WINDOW_WIDTH) {
            return;
        }

        let wrapperRect = document.getElementById(EDITOR_WRAPPER_ID).getBoundingClientRect();
        let {y, height} = this.tableDOM.getBoundingClientRect();
        let menuY = y - wrapperRect.y;

        if(y < 100) {
            menuY -= y-100;
        }

        if(y + height < 100) {
            menuY += y+height-100;
        }

        this.dom.style.top = menuY + "px";
    }

    // close menu on deselect or click outside content area
    mouseUp(e) {
        if(this.tableDOM && !this.dom.contains(e.target) && !this.tableDOM.contains(e.target)) {
            this.close();
        }
        else {
            this.update(this.view);
        }
    }

    close() {
        // reset to default on close
        this.switchState("default");
        this.dom.style.display = "none";
        this.tableDOM = null;
    }

    getDOM() {
        return this.dom;
    }
}

export function tableMenuPlugin(items) {
    let menu = new TableMenu();
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
