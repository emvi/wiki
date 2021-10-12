// Returns the value in an object for given property path.
// The path indexes must be separated by dots, e.g.: my.prop.path
// If the property cannot be found this function returns undefined.
export function findPropValue(obj, prop) {
    if(obj === null || obj === undefined) {
        return obj;
    }

    let path = prop.split('.');
    let current = obj;

    for(let i = 0; i < path.length; i++) {
        if(current[path[i]] === undefined) {
            return undefined;
        }
        else {
            current = current[path[i]];
        }
    }

    return current;
}

export function isEmptyObject(obj) {
    return Object.keys(obj).length === 0 && obj.constructor === Object;
}
