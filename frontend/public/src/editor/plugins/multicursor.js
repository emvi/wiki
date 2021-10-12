import {Plugin} from "prosemirror-state";
import {Decoration, DecorationSet} from "prosemirror-view";

class MultiCursor {
    constructor(updateCursor) {
        this.updateCursor = updateCursor;
        this.decorators = DecorationSet.empty;
        this.lastFrom = -1;
        this.lastTo = -1;

        this.clickEventHandler = () => {this.blur();};
        window.addEventListener("click", this.clickEventHandler);
    }

    destroy() {
        // FIXME gets called somehow...
        //window.removeEventListener("click", this.clickEventHandler)
    }

    init(view) {
        this.view = view;
    }

    // updates the cursor for this user
    updateUserCursor(from, to) {
        if(from !== this.lastFrom || to !== this.lastTo) {
            this.lastFrom = from;
            this.lastTo = to;
            this.updateCursor(from, to);
        }
    }

    // updates the cursor/selection for a remote user
    setCursor(user_id, color, from, to) {
        // create or update cursor
        let cursorDecorators = this.decorators.find(null, null, spec => spec.user_id === user_id && spec.is_cursor);

        if(cursorDecorators && cursorDecorators.length) {
            this.updateCursorPos(cursorDecorators, to, user_id, color);
        }
        else {
            this.createCursor(to, user_id, color);
        }

        // create or update selection
        let selectionDecorators = this.decorators.find(null, null, spec => spec.user_id === user_id && spec.is_selection);
        this.decorators = this.decorators.remove(selectionDecorators);

        if(from !== to) {
            this.createSelection(from, to, user_id, color);
        }
        else if(selectionDecorators.length) {
            this.view.dispatch(this.view.state.tr);
        }
    }

    updateCursorPos(decorators, pos, user_id, color) {
        this.decorators = this.decorators.remove(decorators);
        this.createCursor(pos, user_id, color);
    }

    createCursor(pos, user_id, color) {
        let tr = this.view.state.tr;
        tr.setMeta("update_cursor", true);

        let toDOM = () => {
            let dom = document.createElement("template");
            dom.innerHTML = `<span class="collab-cursor ${color}"></span>`;
            return dom.content.firstChild;
        };

        let decorator = Decoration.widget(pos, toDOM, {user_id, is_cursor: true});
        this.decorators = this.decorators.add(tr.doc, [decorator]);
        this.view.dispatch(tr);
    }

    removeCursor(user_id) {
        let decorators = this.decorators.find(null, null, spec => spec.user_id === user_id);
        this.decorators = this.decorators.remove(decorators);
        this.view.dispatch(this.view.state.tr);
    }

    createSelection(from, to, user_id, color) {
        let tr = this.view.state.tr;
        tr.setMeta("update_cursor", true);
        let decorator = Decoration.inline(from, to, {class: `collab-select ${color}`}, {user_id, is_selection: true});
        this.decorators = this.decorators.add(tr.doc, [decorator]);
        this.view.dispatch(tr);
    }

    getDecorators() {
        return this.decorators;
    }

    mapDecorators(tr) {
        this.decorators = this.decorators.map(tr.mapping, tr.doc);
    }

    blur() {
        if(!this.view.hasFocus()) {
            this.updateUserCursor(-1, -1);
        }
    }
}

export function multicursorPlugin(updateCursor) {
    let multicursor = new MultiCursor(updateCursor);
    let plugin = new Plugin({
        view(view) {
            multicursor.init(view);
            this.multicursor = multicursor;
            return multicursor;
        },
        state: {
            init() {
                return multicursor;
            },
            apply(tr, multicursor) {
                if(tr.docChanged) {
                    multicursor.mapDecorators(tr);
                }
                else if(!tr.getMeta("update_cursor")) {
                    multicursor.updateUserCursor(tr.selection.from, tr.selection.to);
                }

                return multicursor;
            }
        },
        props: {
            decorations(state) {
                let pluginState = this.getState(state);
                return pluginState.getDecorators();
            }
        }
    });

    plugin.setCursor = function(user_id, color, from, to) {
        multicursor.setCursor(user_id, color, from, to);
    };

    plugin.removeCursor = function(user_id) {
        multicursor.removeCursor(user_id);
    };

    plugin.destroy = function() {
        multicursor.destroy();
    };

    return plugin;
}
