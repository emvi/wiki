import CodeMirror from "codemirror";

export function loadCodeMirror(content) {
    setTimeout(() => {
        // TODO make sure this does not loop forever
        if(!content) {
            loadCodeMirror();
            return;
        }

        let codeblocks = [];

        // read all nodes into an array since replacing the nodes with a Codemirror instance also changes the collection.
        for(let codeblock of content.getElementsByTagName("pre")) {
            if(!codeblock.className) {
                codeblocks.push(codeblock);
            }
        }

        for(let i = 0; i < codeblocks.length; i++) {
            let codeblock = codeblocks[i];
            let mimeType = codeblock.firstChild.getAttribute("language");
            let mode = CodeMirror.findModeByMIME(mimeType).mode;

            CodeMirror(editor => {
                codeblock.parentNode.replaceChild(editor, codeblock);
            }, {
                value: codeblock.innerText,
                mode,
                readOnly: true,
                lineNumbers: true,
                indentUnit: 4
            });
        }
    }, 20);
}
