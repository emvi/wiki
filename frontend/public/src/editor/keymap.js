import {keymap} from "prosemirror-keymap";
import {undo, redo} from "prosemirror-history";
import {Selection} from "prosemirror-state";
import {chainCommands, exitCode, toggleMark, joinUp, joinDown} from "prosemirror-commands";
import {undoInputRule} from "prosemirror-inputrules";
import {liftListItem, sinkListItem, splitListItem} from "prosemirror-schema-list";
import {goToNextCell} from "prosemirror-tables";

export function keymapPlugin(schema) {
    let brCmd = chainCommands(exitCode, (state, dispatch) => {
        if(dispatch) {
            dispatch(state.tr.replaceSelectionWith(schema.nodes.hard_break.create()).scrollIntoView());
        }

        return true;
    });

    return keymap({
        // text formatting
        "Mod-b": toggleMark(schema.marks.bold),
        "Mod-i": toggleMark(schema.marks.italic),
        "Mod-u": toggleMark(schema.marks.underlined),

        // enter
        "Enter": chainCommands(enterHandler(schema.nodes.paragraph), splitListItem(schema.nodes.list_item), splitListItem(schema.nodes.check_list_item)),

        // lists and tables
        "Tab": chainCommands(sinkListItem(schema.nodes.list_item), sinkListItem(schema.nodes.check_list_item), goToNextCell(1)),
        "Shift-Tab": chainCommands(liftListItem(schema.nodes.list_item), liftListItem(schema.nodes.check_list_item), goToNextCell(-1)),

        // history
        "Ctrl-z": undo,
        "Ctrl-y": redo,

        // insert line break
        "Mod-Enter": brCmd,
        "Shift-Enter": brCmd,

        // merge blocks
        "Alt-ArrowUp": joinUp,
        "Alt-ArrowDown": joinDown,

        // general keys
        "Backspace": undoInputRule,

        // code editor
        "ArrowLeft": arrowHandler("left"),
        "ArrowRight": arrowHandler("right"),
        "ArrowUp": arrowHandler("up"),
        "ArrowDown": arrowHandler("down")
    });
}

function enterHandler(paragraph) {
    return (state, dispatch) => {
        let $from = state.selection.$from;

        if($from.depth > 1 && $from.node(-1).type.name === "image") {
            if(dispatch) {
                let pos = state.doc.resolve($from.after($from.depth-1));
                let tr = state.tr;

                // if there is no node after the image, create an empty paragraph
                // else jump into paragraph
                if(!pos.nodeAfter) {
                    tr.setSelection(new Selection(pos, pos));
                    tr.replaceSelectionWith(paragraph.create());
                }
                else {
                    tr.setSelection(Selection.near(pos));
                }

                dispatch(tr);
            }

            return true;
        }

        return false;
    }
}

// used to handle arrow keys in code editor
function arrowHandler(dir) {
    return (state, dispatch, view) => {
        if(state.selection.empty && view.endOfTextblock(dir)) {
            let side = dir === "left" || dir === "up" ? -1 : 1;
            let $head = state.selection.$head;

            if($head.depth) {
                let nextPos = Selection.near(state.doc.resolve(side > 0 ? $head.after() : $head.before()), side);

                if(nextPos.$head && nextPos.$head.parent.type.name === "code_block") {
                    dispatch(state.tr.setSelection(nextPos));
                    return true;
                }
            }
        }

        return false;
    }
}
