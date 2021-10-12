export function tagsToStringArray(tags) {
    let tagStr = [];

    for(let i = 0; i < tags.length; i++){
        tagStr.push(tags[i].name);
    }

    return tagStr;
}

export function getTextWidth(text, font) {
    // re-use canvas object for better performance
    let canvas = getTextWidth.canvas || (getTextWidth.canvas = document.createElement("canvas"));
    let context = canvas.getContext("2d");
    context.font = font;
    let metrics = context.measureText(text);
    return Math.ceil(metrics.width);
}
