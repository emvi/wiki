const INFOBOX_COLORS = [
    "blue",
    "green",
    "orange",
    "red"
];

export class InfoboxView {
    constructor(node, view, getPos) {
        this.node = node;
        this.view = view;
        this.getPos = getPos;
        this.buildContentDOM();
        this.buildDOM();
    }

    buildContentDOM() {
        let dom = document.createElement("template");
        dom.innerHTML = `<span class="infobox-content"></span>`;
        this.contentDOM = dom.content.firstChild;
    }

    buildDOM() {
        let color = this.node.attrs.color;
        let dom = document.createElement("template");
        dom.innerHTML = `<div class="infobox ${color}" color="${color}"></div>`;
        let node = dom.content.firstChild;
        node.appendChild(this.contentDOM);
        node.appendChild(this.buildColorSelectDOM());
        this.dom = node;
        this.setColor(color);
    }

    buildColorSelectDOM() {
        let dom = document.createElement("template");
        dom.innerHTML = `<div class="infobox-colors" contenteditable="false"></div>`;
        let node = dom.content.firstChild;
        this.colorNodes = [];

        for (let i = 0; i < INFOBOX_COLORS.length; i++) {
            let c = INFOBOX_COLORS[i];
            let option = document.createElement("template");
            option.innerHTML = `<div class="infobox-color bg-${c}-100"></div>`;
            let colorNode = option.content.firstChild;

            colorNode.addEventListener("click", () => {
                this.setColor(c);
                this.setNodeAttribute(c);
            });

            node.appendChild(colorNode);
            this.colorNodes.push(colorNode);
        }

        this.colorNode = node;
        return node;
    }

    setNodeAttribute(color) {
        let tr = this.view.state.tr;
        tr.setNodeMarkup(this.getPos(), null, {color});
        this.view.dispatch(tr);
    }

    setColor(color) {
        for(let i = 0; i < INFOBOX_COLORS.length; i++) {
            this.dom.classList.remove(`${INFOBOX_COLORS[i]}`);
            this.colorNodes[i].classList.remove("active");

            if(INFOBOX_COLORS[i] === color) {
                this.colorNodes[i].classList.add("active");
            }
        }

        this.dom.classList.add(`${color}`);
        this.activeColor = color;
    }

    update(node) {
        if(node.type.name === "infobox") {
            this.node = node;

            if (this.activeColor !== node.attrs.color) {
                this.setColor(node.attrs.color);
            }

            return true;
        }

        return false;
    }

    stopEvent(e) {
        return e.target && this.colorNode.contains(e.target);
    }

    ignoreMutation(e) {
        return e.target && this.colorNode.contains(e.target);
    }
}
