export function addTableEnlargeHandles(content) {
    let dom = document.createElement("div");
    dom.innerHTML = content;
    let tableWrappers = findElementsByTagNameAndClassName(dom, "DIV", "table-wrapper");

    for(let wrapper of tableWrappers) {
        let handle = document.createElement("div");
        handle.classList.add("table-handle", "icon", "icon-enlarge", "size-40", "cursor-pointer");
        wrapper.insertBefore(handle, wrapper.firstChild);
    }

    return dom.innerHTML;
}

function findElementsByTagNameAndClassName(node, tagName, className) {
    let results = [];

    if(node.tagName === tagName && node.classList.contains(className)) {
        results.push(node);
    }

    for(let child of node.childNodes) {
        if(child.childNodes && child.childNodes.length) {
            results = results.concat(findElementsByTagNameAndClassName(child, tagName, className));
        }
    }

    return results;
}
