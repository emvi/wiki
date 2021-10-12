import {setBlockType, toggleMark} from "prosemirror-commands";
import {wrapInList} from "prosemirror-schema-list";
import {undo, redo} from "prosemirror-history";
import {wrapIn} from "prosemirror-commands";
import {
    addColumnAfter,
    addColumnBefore,
    addRowAfter,
    addRowBefore,
    deleteColumn,
    deleteRow, deleteTable, mergeCells, setCellAttr, splitCell, toggleHeaderCell,
    toggleHeaderColumn, toggleHeaderRow
} from "prosemirror-tables";

export function menuCommands(editorSchema, editor) {
    return [
        {command: undo, action: {tooltip: "Undo", icon: "icon-undo"}},
        {command: redo, action: {tooltip: "Redo", icon: "icon-redo"}},
        {command: insertText(editorSchema, "@"), action: {tooltip: "Mention", icon: "icon-at"}},
        {command: selectFile(editorSchema.nodes.image, editor.uploadPlugin, true), action: {tooltip: "Upload image", icon: "icon-img"}},
        {command: selectFile(editorSchema.nodes.file, editor.uploadPlugin, false), action: {tooltip: "Upload file", icon: "icon-file"}},
        {command: wrapInList(editorSchema.nodes.bullet_list), action: {tooltip: "Unordered list", icon: "icon-ul"}},
        {command: wrapInList(editorSchema.nodes.ordered_list), action: {tooltip: "Ordered list", icon: "icon-ol"}},
        {command: wrapInList(editorSchema.nodes.check_list), action: {tooltip: "Check list", icon: "icon-checkbox"}},
        {command: insertTable(editorSchema.nodes.table, editorSchema.nodes.table_row, editorSchema.nodes.table_cell, editorSchema.nodes.paragraph), action: {tooltip: "Table", icon: "icon-table"}},
        {command: wrapIn(editorSchema.nodes.infobox), action: {tooltip: "Infobox", icon: "icon-infobox"}},
        {command: wrapIn(editorSchema.nodes.blockquote), action: {tooltip: "Blockquote", icon: "icon-quote"}},
        {command: setBlockType(editorSchema.nodes.code_block), action: {tooltip: "Code", icon: "icon-code"}},
        {command: insertHr(editorSchema.nodes.horizontal_rule), action: {tooltip: "Horizontal rule", icon: "icon-hr"}}
    ];
}

export function inlineMenuCommands(editorSchema) {
    return [
        {command: setBlockType(editorSchema.nodes.paragraph), action: {tooltip: "Plain text", icon: "icon-p"}},
        {command: setBlockType(editorSchema.nodes.headline, {level: 2}), action: {tooltip: "Headline 1", icon: "icon-h1"}},
        {command: setBlockType(editorSchema.nodes.headline, {level: 3}), action: {tooltip: "Headline 2", icon: "icon-h2"}},
        {command: setBlockType(editorSchema.nodes.headline, {level: 4}), action: {tooltip: "Headline 3", icon: "icon-h3"}},
        {command: toggleMark(editorSchema.marks.bold), action: {tooltip: "Bold", icon: "icon-bold"}},
        {command: toggleMark(editorSchema.marks.italic), action: {tooltip: "Italic", icon: "icon-italic"}},
        {command: toggleMark(editorSchema.marks.strikethrough), action: {tooltip: "Strikethrough", icon: "icon-strike"}},
        {command: toggleMark(editorSchema.marks.sup), action: {tooltip: "Superscript", icon: "icon-sup"}},
        {command: toggleMark(editorSchema.marks.sub), action: {tooltip: "Subscript", icon: "icon-sub"}},
        {command: toggleMark(editorSchema.marks.underlined), action: {tooltip: "Highlight", icon: "icon-u"}},
        {command: insertLink(editorSchema), state: "insertlink", action: {tooltip: "Link", icon: "icon-link"}},
        {command: wrapInList(editorSchema.nodes.bullet_list), action: {tooltip: "Unordered list", icon: "icon-ul"}},
        {command: wrapInList(editorSchema.nodes.ordered_list), action: {tooltip: "Ordered list", icon: "icon-ol"}},
        {command: wrapInList(editorSchema.nodes.check_list), action: {tooltip: "Check list", icon: "icon-checkbox"}},
        {command: wrapIn(editorSchema.nodes.infobox), action: {tooltip: "Infobox", icon: "icon-infobox"}},
        {command: wrapIn(editorSchema.nodes.blockquote), action: {tooltip: "Blockquote", icon: "icon-quote"}},
        {command: toggleMark(editorSchema.marks.code), action: {tooltip: "Code", icon: "icon-code"}}
    ];
}

export function tableMenuCommands() {
    return [
        {command: setCellAttr, state: "color", action: {tooltip: "Cell color", icon: "icon-color"}},
        {command: toggleHeaderColumn, action: {tooltip: "Toggle header column", icon: "icon-header-column"}},
        {command: addColumnBefore, action: {tooltip: "Insert column before", icon: "icon-add-column icon-rotate-180"}},
        {command: addColumnAfter, action: {tooltip: "Insert column after", icon: "icon-add-column"}},
        {command: deleteColumn, action: {tooltip: "Delete column", icon: "icon-delete-column"}},
        {command: toggleHeaderRow, action: {tooltip: "Toggle header row", icon: "icon-header-row"}},
        {command: addRowBefore, action: {tooltip: "Insert row before", icon: "icon-add-row icon-rotate-180"}},
        {command: addRowAfter, action: {tooltip: "Insert row after", icon: "icon-add-row"}},
        {command: deleteRow, action: {tooltip: "Delete row", icon: "icon-delete-row"}},
        {command: toggleHeaderCell, action: {tooltip: "Toggle header cell", icon: "icon-header-cell"}},
        {command: mergeCells, action: {tooltip: "Merge cells", icon: "icon-merge-cell"}},
        {command: splitCell, action: {tooltip: "Split cell", icon: "icon-split-cell"}},
        {command: deleteTable, action: {tooltip: "Delete table", icon: "icon-delete-table"}}
    ];
}

function insertText(schema, text) {
    return function(state, dispatch) {
        if(dispatch) {
            dispatch(state.tr.replaceSelectionWith(schema.text(text)).scrollIntoView());
        }

        return true;
    }
}

function insertHr(hr) {
    return function(state, dispatch) {
        let insertable = canInsert(state, hr);

        if(!dispatch) {
            return insertable;
        }

        if(dispatch && insertable) {
            dispatch(state.tr.replaceSelectionWith(hr.create()));
        }
    }
}

function insertTable(table, row, cell, paragraph) {
    return function(state, dispatch) {
        let insertable = canInsert(state, table);

        if(!dispatch) {
            return insertable;
        }

        if(dispatch && insertable) {
            // create empty 2x2 table
            let td = cell.create({colwidth: [100]}, paragraph.create());
            let tr = row.create(null, [td, td]);
            let node = table.create(null, [tr, tr]);
            dispatch(state.tr.replaceSelectionWith(node));
        }
    }
}

function selectFile(node, uploadPlugin, imageOnly) {
    return function(state, dispatch) {
        let insertable = canInsert(state, node);

        if(!dispatch) {
            return insertable;
        }

        if(dispatch && insertable) {
            let input = document.createElement("input");
            input.type = "file";

            if(imageOnly) {
                input.accept = "image/*";
            }

            input.addEventListener("change", e => {
                e.stopPropagation();
                e.preventDefault();
                let pos = state.tr.selection.from;
                uploadPlugin.spec.fileUploader.uploadFile(e.target.files[0], pos);
                return false;
            });

            input.click();
        }
    }
}

function insertLink(editorSchema) {
    let markType = editorSchema.marks.link;
    let canInsert = toggleMark(markType);

    return function (state, dispatch, attrs) {
        let insertable = canInsert(state);

        if(!dispatch) {
            return insertable;
        }

        if(dispatch && insertable) {
            let active = markActive(state, markType);

            // remove if active and input is empty
            if(attrs.href.length === 0 && active) {
                toggleMark(editorSchema.marks.link)(state, dispatch);
            }
            // update
            else if(attrs.href.length !== 0) {
                let tr = state.tr;
                let {from, to} = state.selection;

                if(state.doc.rangeHasMark(from, to, markType)) {
                    tr.removeMark(from, to, markType);
                }

                dispatch(tr.addMark(from, to, markType.create(attrs)));
            }
        }
    }
}

function canInsert(state, nodeType) {
    let $from = state.selection.$from;

    for (let d = $from.depth; d >= 0; d--) {
        let index = $from.index(d);

        if($from.node(d).canReplaceWith(index, index, nodeType)) {
            return true;
        }
    }

    return false;
}

function markActive(state, type) {
    let {from, $from, to, empty} = state.selection;

    if(empty) {
        return type.isInSet(state.storedMarks || $from.marks());
    }

    return state.doc.rangeHasMark(from, to, type);
}
