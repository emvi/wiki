// Returns true if the given parameter is an empty object.
export function isEmptyObject(obj) {
    return isObject(obj) && Object.keys(obj).length === 0 && obj.constructor === Object;
}

// Returns true if the given parameter is an object.
export function isObject(obj) {
    return typeof obj === "object" && obj !== null && obj !== undefined;
}
