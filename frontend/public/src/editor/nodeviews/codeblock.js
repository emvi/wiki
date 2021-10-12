import CodeMirror from "codemirror";
import "codemirror/mode/meta.js";
import {exitCode} from "prosemirror-commands";
import {undo, redo} from "prosemirror-history";
import {TextSelection, Selection} from "prosemirror-state";
import {editorSchema as schema} from "../schema.js";
import {editorSchema} from "../schema";

// import all codemirror languages
var req = require.context("../../../node_modules/codemirror/mode", true, /^.*\.js$/);

req.keys().forEach(key => {
    req(key);
});

export class CodeBlockView {
    constructor(node, view, getPos) {
        this.node = node;
        this.view = view;
        this.getPos = getPos;
        this.incomingChanges = false;
        this.updating = false;

        this.cm = new CodeMirror(null, {
            value: this.node.textContent,
            lineNumbers: true,
            extraKeys: this.codeMirrorKeymap(),
            indentUnit: 4
        });
        let dom = this.cm.getWrapperElement();
        this.addLanguageSelection(dom);
        this.dom = dom;
        setTimeout(() => {this.cm.refresh();}, 20);

        this.cm.on("beforeChange", () => {
            this.incomingChanges = true;
        });
        this.cm.on("cursorActivity", () => {
            if(!this.updating && !this.incomingChanges) {
                this.forwardSelection();
            }
        });
        this.cm.on("changes", () => {
            if(!this.updating) {
                this.valueChanged();
                this.forwardSelection();
            }

            this.incomingChanges = false
        });
        this.cm.on("focus", () => {
            this.forwardSelection();
        });
    }

    addLanguageSelection(dom) {
        let options = [];

        for(let i = 0; i < CodeMirror.modeInfo.length; i++) {
            options.push(`<option value="${CodeMirror.modeInfo[i].mime}">${CodeMirror.modeInfo[i].name}</option>`);
        }

        this.langSelect = document.createElement("select");
        this.langSelect.className = "code-editor-lang-select";
        this.langSelect.innerHTML = options.join();
        this.langSelect.addEventListener("change", () => {
            this.setNodeLanguageAttribute(this.langSelect.value);
            this.setMode(this.langSelect.value);
        });
        dom.appendChild(this.langSelect);

        // initial language selection
        this.langSelect.value = this.node.attrs.language;
        this.setMode(this.langSelect.value);
    }

    setNodeLanguageAttribute(mimeType) {
        let tr = this.view.state.tr;
        tr.setNodeMarkup(this.getPos(), null, {language: mimeType});
        this.view.dispatch(tr);
    }

    setMode(mimeType) {
        let mode = CodeMirror.findModeByMIME(mimeType).mode;
        this.cm.setOption("mode", mode);
    }

    forwardSelection() {
        if(!this.cm.hasFocus()) {
            return;
        }

        let state = this.view.state;
        let selection = this.asProseMirrorSelection(state.doc);

        if(!selection.eq(state.selection)) {
            this.view.dispatch(state.tr.setSelection(selection));
        }
    }

    asProseMirrorSelection(doc) {
        let offset = this.getPos() + 1;
        let anchor = this.cm.indexFromPos(this.cm.getCursor("anchor")) + offset;
        let head = this.cm.indexFromPos(this.cm.getCursor("head")) + offset;
        return TextSelection.create(doc, anchor, head);
    }

    setSelection(anchor, head) {
        this.cm.focus();
        this.updating = true;
        this.cm.setSelection(this.cm.posFromIndex(anchor), this.cm.posFromIndex(head));
        this.updating = false;
    }

    valueChanged() {
        let change = computeChange(this.node.textContent, this.cm.getValue());

        if(change) {
            let start = this.getPos() + 1;
            let tr = this.view.state.tr.replaceWith(start + change.from,
                start + change.to,
                change.text ? schema.text(change.text) : null);
            this.view.dispatch(tr);
        }
    }

    codeMirrorKeymap() {
        let view = this.view;
        let mod = /Mac/.test(navigator.platform) ? "Cmd" : "Ctrl";

        return CodeMirror.normalizeKeyMap({
            Up: () => {
                return this.maybeEscape("line", -1);
            },
            Left: () => {
                return this.maybeEscape("char", -1);
            },
            Down: () => {
                return this.maybeEscape("line", 1);
            },
            Right: () => {
                return this.maybeEscape("char", 1);
            },
            [`${mod}-Z`]: () => {
                return undo(view.state, view.dispatch);
            },
            [`Shift-${mod}-Z`]: () => {
                return redo(view.state, view.dispatch);
            },
            [`${mod}-Y`]: () => {
                return redo(view.state, view.dispatch);
            },
            "Ctrl-Enter": () => {
                if(exitCode(view.state, view.dispatch)) {
                    view.focus();
                }
            },
            "Backspace": () => {
                let pos = this.cm.getCursor();

                if(pos.line === 0 && pos.ch === 0) {
                    this.deleteCodeBlock();
                    return;
                }

                return CodeMirror.Pass;
            }
        })
    }

    deleteCodeBlock() {
        // range doesn't work if the editor is the first node
        let pos = this.getPos();
        let tr = this.view.state.tr;
        tr.setBlockType(pos, pos+1, editorSchema.nodes.paragraph);
        this.view.dispatch(tr);
        this.view.focus();
    }

    maybeEscape(unit, dir) {
        let pos = this.cm.getCursor();

        if(this.cm.somethingSelected() ||
            pos.line !== (dir < 0 ? this.cm.firstLine() : this.cm.lastLine()) ||
            (unit === "char" && pos.ch !== (dir < 0 ? 0 : this.cm.getLine(pos.line).length))) {
            return CodeMirror.Pass;
        }

        this.view.focus();
        let targetPos = this.getPos() + (dir < 0 ? 0 : this.node.nodeSize);
        let selection = Selection.near(this.view.state.doc.resolve(targetPos), dir);
        this.view.dispatch(this.view.state.tr.setSelection(selection).scrollIntoView());
        this.view.focus();
    }

    update(node) {
        if(node.type !== this.node.type) {
            return false;
        }

        if(this.node.attrs.language !== node.attrs.language) {
            this.langSelect.value = node.attrs.language;
            this.setMode(node.attrs.language);
        }

        this.node = node;
        let change = computeChange(this.cm.getValue(), node.textContent);

        if(change) {
            this.updating = true;
            this.cm.replaceRange(change.text, this.cm.posFromIndex(change.from), this.cm.posFromIndex(change.to));
            this.updating = false;
        }

        return true;
    }

    selectNode() {
        this.cm.focus();
    }

    stopEvent() {
        return true;
    }
}

function computeChange(oldVal, newVal) {
    if(oldVal === newVal) {
        return null;
    }

    let start = 0, oldEnd = oldVal.length, newEnd = newVal.length;

    while(start < oldEnd && oldVal.charCodeAt(start) === newVal.charCodeAt(start)) {
        ++start;
    }

    while(oldEnd > start && newEnd > start && oldVal.charCodeAt(oldEnd - 1) === newVal.charCodeAt(newEnd - 1)) {
        oldEnd--;
        newEnd--;
    }

    return {from: start, to: oldEnd, text: newVal.slice(start, newEnd)};
}
