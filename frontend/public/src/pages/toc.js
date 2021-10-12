import {scrollTo} from "../util";
import {buildTableOfContents} from "../editor";

const tocPadding = 100;

export const TocMixin = {
    data() {
        return {
            toc: []
        };
    },
    mounted() {
        window.addEventListener("scroll", this.scroll);
    },
    beforeDestroy() {
        window.removeEventListener("scroll", this.scroll);
    },
    methods: {
        buildTableOfContents(content) {
            let dom = document.createElement("div");
            dom.innerHTML = content;
            this.toc = buildTableOfContents(dom);
            this.scroll();
        },
        jumpToc(headline) {
            let element = this._getContent().getElementsByTagName(headline.tag)[headline.number-1];
            scrollTo(element, tocPadding);
        },
        scroll() {
            if(this.toc.length) {
                let headlines = [];
                this._flattenList(this.toc, headlines);
                let active = this._findLastHeadline(headlines, tocPadding);
                headlines[active].active = true;
                this._setParentHeadlinesActive(headlines, active);
            }
        },
        _getContent() {
            let content = this.$refs.content;

            if(!content) {
                return document.createElement("div");
            }

            return content;
        },
        _flattenList(toc, list) {
            for(let i = 0; i < toc.length; i++) {
                list.push(toc[i]);

                if(toc[i].children) {
                    this._flattenList(toc[i].children, list);
                }
            }
        },
        _findLastHeadline(toc, padding) {
            let content = this._getContent();
            let active = 0;

            for(let i = 0; i < toc.length; i++) {
                let headline = toc[i];
                let element = content.getElementsByTagName(headline.tag)[headline.number - 1];

                if (element) {
                    let top = element.getBoundingClientRect().top - padding;
                    headline.active = false;

                    if (Math.floor(top) <= 0) {
                        active = i;
                    }
                }
            }

            return active;
        },
        _setParentHeadlinesActive(headlines, active) {
            if(headlines[active].tag === "H2") {
                return;
            }

            active--;

            while(active > -1 && headlines[active].tag !== "H2") {
                active--;
            }

            if(active > -1) {
                headlines[active].active = true;
            }
        }
    }
}
