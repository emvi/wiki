import Mark from "mark.js";

// getQueryFromCmd returns the query part from a command.
// Example: .page -> page
export function getQueryFromCmd(cmd) {
    return cmd.substring(1).trim().toLowerCase();
}

// nameStartsWith returns true if one of the names starts with given query.
// This is used to match search queries against aliases.
export function nameStartsWith(names, query) {
    for(let i = 0; i < names.length; i++) {
        if(names[i].startsWith(query)) {
            return true;
        }
    }

    return false;
}

// markInText highlights the search query keywords in text and returns the title as HTML.
export function markInText(query, text) {
    let dom = document.createElement("div");
    dom.innerText = text;
    let mark = new Mark(dom);
    mark.mark(query.split(" "), {element: "u"});
    return dom.innerHTML;
}

// updateSelectedRow updates the selected row depending on the number of result in store.
export function updateSelectedRow(row, results, store) {
    if(!results) {
        return;
    }

    if(row < 0) {
        store.dispatch("selectRow", results-1);
    }
    else if(row > results-1) {
        store.dispatch("selectRow", 0);
    }
}

// updateSelection returns valid index based on maxIndex.
// index > maxIndex -> 0
// index < 0 -> maxIndex
export function updateSelection(index, maxIndex) {
    if(index > maxIndex) {
        return 0;
    }

    if(index < 0) {
        return maxIndex;
    }

    return index;
}
