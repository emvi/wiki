import {Plugin} from "prosemirror-state";
import {editorSchema} from "../schema.js";

const EDITOR_WRAPPER_ID = "editor-wrapper";
const MAX_WINDOW_WIDTH = 512;

class ItemMenu {
    constructor(menu, view, items) {
        this.menu = menu;
        this.view = view;
        this.items = items;
        this.buildDOM();
    }

    buildDOM() {
        this.dom = document.createElement("div");
        this.dom.className = "inline-menu";

        for(let item of this.items) {
            item.dom = ItemMenu.buildItem(item.action.tooltip, item.action.icon);
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

    enableRelocate() {
        return true;
    }

    setCommand(command) {}

    getDOM() {
        return this.dom;
    }
}

class InsertLinkMenu {
    constructor(menu, view) {
        this.menu = menu;
        this.view = view;
        this.buildDOM();
    }

    buildDOM() {
        this.dom = document.createElement("div");
        this.dom.className = "inline-menu inline-link";

        this.input = InsertLinkMenu.buildInput();
        this.handleEnter(this.view, this.input, () => {this.confirm();});
        let confirm = InsertLinkMenu.buildButton("Confirm", "icon-check");
        this.handleClick(this.view, confirm, () => {this.confirm();});
        let cancel = InsertLinkMenu.buildButton("Cancel", "icon-close");
        this.handleClick(this.view, cancel, () => {this.cancel();});

        this.dom.appendChild(this.input);
        this.dom.appendChild(confirm);
        this.dom.appendChild(cancel);
    }

    static buildInput() {
        let input = document.createElement("input");
        input.type = "url";
        input.name = "link";
        input.placeholder = "http://...";
        return input;
    }

    static buildButton(tooltip, icon) {
        let i = document.createElement("i");
        i.className = "icon "+icon;
        i.title = tooltip;
        return i;
    }

    handleEnter(view, input, callback) {
        input.addEventListener("keyup", e => {
            if(e.keyCode === 13) {
                e.preventDefault();
                view.focus();
                callback();
            }
        });
    }

    handleClick(view, button, callback) {
        button.addEventListener("click", e => {
            e.preventDefault();
            view.focus();
            callback();
        });
    }

    confirm() {
        let value = this.input.value.trim();

        if(value.length && !value.toLowerCase().startsWith("http") && !value.toLowerCase().startsWith("mailto:")) {
            value = `http://${value}`;
        }

        this.command(this.view.state, this.view.dispatch, {href: value});
        this.cancel();
    }

    cancel() {
        this.menu.switchState("default");
    }

    enableItems() {
        return true;
    }

    activate() {
        let config = {childList: true};
        let observer = new MutationObserver(mutations => {
            for(let mutation of mutations) {
                if(mutation.addedNodes.length === 1 && mutation.addedNodes[0] === this.dom) {
                    this.initInput();
                }
            }

            observer.disconnect();
        });
        observer.observe(this.view.dom.parentNode, config);
    }

    initInput() {
        let state = this.view.state;
        let {from, to} = state.selection;

        if(state.doc.rangeHasMark(from, to, editorSchema.marks.link)) {
            state.doc.nodesBetween(from, to, node => {
                if(node.isLeaf && node.marks.length) {
                    for(let mark of node.marks) {
                        if(mark.type.name === "link") {
                            this.input.value = mark.attrs.href;
                            return false;
                        }
                    }
                }

                return true;
            });
        }

        this.input.focus();
    }

    reset() {
        this.input.value = "";
    }

    enableRelocate() {
        return false;
    }

    setCommand(command) {
        this.command = command;
    }

    getDOM() {
        return this.dom;
    }
}

export class EditorInlineMenu {
    init(view, items) {
        this.view = view;
        this.buildDOM(items);

        this.cursorInMenu = false;
        this.mousedownEventHandler = e => {this.mouseDown(e);};
        window.addEventListener("mousedown", this.mousedownEventHandler);
        this.mouseupEventHandler = e => {this.mouseUp(e);};
        window.addEventListener("mouseup", this.mouseupEventHandler);
    }

    destroy() {
        // FIXME gets called somehow...
        //window.removeEventListener("mouseup", this.mouseupEventHandler);
    }

    buildDOM(items) {
        this.states = {
            default: new ItemMenu(this, this.view, items),
            insertlink: new InsertLinkMenu(this, this.view)
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
        this.update(this.view, null);
    }

    replaceDOM() {
        let dom = this.activeState.getDOM();
        this.calcPosition(dom);

        if(this.dom) {
            this.view.dom.parentNode.replaceChild(dom, this.dom);
            this.dom = dom;
        }
        else {
            this.dom = dom;
        }
    }

    update(view, lastState) {
        let state = view.state;

        if(lastState && lastState.doc.eq(state.doc) && lastState.selection.eq(state.selection)) {
            return;
        }

        if(state.selection.empty || !this.activeState.enableItems()) {
            this.close();
            return;
        }

        this.calcPosition(this.dom);
    }

    calcPosition(dom) {
        if(!this.activeState.enableRelocate()) {
            return;
        }

        const selection = window.getSelection();

        if(!selection.rangeCount) {
            return;
        }

        if(window.innerWidth >= MAX_WINDOW_WIDTH) {
            const selectionRange = selection.getRangeAt(0);
            let wrapperRect = document.getElementById(EDITOR_WRAPPER_ID).getBoundingClientRect();
            wrapperRect.x += window.scrollX;
            wrapperRect.y += window.scrollY;
            let {x, y, width} = selectionRange.getBoundingClientRect();
            x -= wrapperRect.x;
            y -= wrapperRect.y;
            dom.style.left = x + (width / 2) + "px";
            dom.style.top = y + window.scrollY - 10 + "px";
        }

        dom.style.display = "";
    }

    mouseUp() {
        let focused = this.view.dom.classList.contains("ProseMirror-focused");

        if(!focused && !this.cursorInMenu) {
            this.close();
        }
    }

    mouseDown(e) {
        this.cursorInMenu = this.dom.contains(e.target);
    }

    close() {
        if(this.dom) {
            // reset to default on close
            this.switchState("default");
            this.dom.style.display = "none";
        }
    }

    getDOM() {
        return this.dom;
    }
}

export function inlineMenuPlugin(items) {
    let menu = new EditorInlineMenu();
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
