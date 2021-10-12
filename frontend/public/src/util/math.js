export function getSizeFromBytes(size) {
    if(typeof size !== "number") {
        return size;
    }

    let i = 0;

    while(size > 1000) {
        size /= 1000;
        i++;
    }

    let units = ['bytes', 'kB', 'MB', 'GB'];

    if(i === 0) {
        return `${size} bytes`;
    }

    return `${size.toFixed(2)} ${units[i]}`
}
