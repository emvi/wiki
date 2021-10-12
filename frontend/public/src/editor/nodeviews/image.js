const CAPTION_PLACEHOLDER = "Image caption (optional)";

export class ImageView {
    constructor(node, view, getPos) {
        this.node = node;
        this.view = view;
        this.getPos = getPos;
        this.buildContentDOM();
        this.buildDOM();
        this.update(node);
    }

    buildContentDOM() {
        let dom = document.createElement("template");
        dom.innerHTML = `<span></span>`;
        this.contentDOM = dom.content.firstChild;
    }

    buildDOM() {
        let src = this.node.attrs.src;
        let dom = document.createElement("template");
        dom.innerHTML = `<figure src="${src}"></figure>`;
        let node = dom.content.firstChild;
        node.appendChild(this.buildImgDOM());
        node.appendChild(this.buildCaptionDOM());
        this.dom = node;
    }

    buildImgDOM() {
        let img = document.createElement("img");
        img.src = this.node.attrs.src;
        return img;
    }

    buildCaptionDOM() {
        let dom = document.createElement("template");
        dom.innerHTML = `<figcaption><div class="placeholder"><p></p></div></figcaption>`;
        let node = dom.content.firstChild;
        node.appendChild(this.contentDOM);
        this.captionDOM = node;
        return node;
    }

    update(node) {
        if(node.type.name === "image") {
            this.togglePlaceholder(node);
            return true;
        }

        return false;
    }

    togglePlaceholder(node) {
        if(node.firstChild.content.size === 0) {
            this.showPlaceholder();
        }
        else {
            this.hidePlaceholder();
        }
    }

    showPlaceholder() {
        this.captionDOM.firstChild.firstChild.innerText = CAPTION_PLACEHOLDER;
    }

    hidePlaceholder() {
        this.captionDOM.firstChild.firstChild.innerText = "";
    }
}
