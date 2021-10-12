export function findIndexById(list, id) {
    if(!list || !list.length) {
        return -1;
    }

    for(let i = 0; i < list.length; i++) {
        if(list[i].id === id) {
            return i;
        }
    }

    return -1;
}

export function addAttrToListElements(list, attr, value) {
    for(let i = 0; i < list.length; i++) {
        list[i][attr] = value;
    }

    return list;
}

export function removeFromListByIds(list, existing) {
    let existingMap = new Map();

    for(let i = 0; i < existing.length; i++) {
        existingMap.set(existing[i].id, true);
    }

    let out = [];

    for(let i = 0; i < list.length; i++) {
        if(!existingMap.get(list[i].id)) {
            out.push(list[i]);
        }
    }

    return out;
}
