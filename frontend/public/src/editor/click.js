// Handles the click on elements inside the editor.
export function handleClickOn(editorView, pos, node, nodePos, event) {
    // refocus editor on placeholder click
    // this is a workaround for a FireFox bug
    if(event.target && event.target.children.length > 0 && event.target.children[0].classList.contains("editor-placeholder")) {
        editorView.focus();
    }

    if(event.target.classList.contains("checklist-trigger")) {
        let tr = editorView.state.tr.setNodeMarkup(nodePos, null, {checked: !node.attrs.checked});
        editorView.dispatch(tr);
        return true;
    }
}
