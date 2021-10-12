import {Plugin} from "prosemirror-state";
import {Decoration, DecorationSet} from "prosemirror-view";
import {SearchService} from "../../service";
import {debounce} from "../../util";
import {editorSchema} from "../schema";

const EDITOR_WRAPPER_ID = "editor-wrapper";
const SEARCH_PLACEHOLDER = `<li class="placeholder">Type to search...</li>`;
const mentionRegex = new RegExp("(^|\\s|[^\\wöäüÖÄÜß])@([\\w-+öäüÖÄÜß]*\\s?[\\w-+öäüÖÄÜß]*)$");

class Mentions {
    constructor() {
        this.search = debounce(() => {
            if(!this.match) {
                return;
            }

            let filter = {
                articles: true,
                lists: true,
                user: true,
                groups: true,
                tags: true,
                feed: false,
                articles_limit: 3,
                lists_limit: 3,
                user_limit: 5,
                groups_limit: 3,
                tags_limit: 3
            };

            SearchService.findAll(this.match.query, filter, 5)
                .then(results => {
                    this.selectedResultIndex = 0;
                    this.updateResultsDOM(results);
                });
        }, 200);
    }

    init(view) {
        this.view = view;
        this.match = null;
        this.results = [];
        this.selectedResultIndex = 0;
        this.active = false;
        this.buildDOM();
    }

    destroy() {}

    buildDOM() {
        this.dom = document.createElement("ul");
        this.dom.className = "mention-overlay";
        this.dom.style.display = "none";
    }

    update(tr) {
        const $position = this.getPosition(tr);

        if(!$position) {
            return;
        }

        this.match = this.getMatch(tr, $position);

        if(!this.match) {
            this.close();
        }
        else {
            this.active = true;
            this.search();
            this.updateDOM();
        }
    }

    getPosition(tr) {
        let selection = tr.selection;

        if(!selection || selection.from !== selection.to) {
            return null;
        }

        return selection.$from;
    }

    getMatch(tr, $position) {
        let textFrom = $position.depth === 0 ? $position.pos : $position.before();
        let text = tr.doc.textBetween(textFrom, $position.pos, "\n", "\0");
        let match = text.match(mentionRegex);

        if(!match) {
            return null;
        }

        let from = match.index+match[1].length;
        from += $position.start();
        let to = from+match[0].length;

        if(match[1].length) {
            to--;
        }

        return {
            from,
            to,
            query: match[2].trim()
        };
    }

    close() {
        this.match = null;
        this.results = [];
        this.selectedResultIndex = 0;
        this.active = false;
        this.dom.style.display = "none";
    }

    updateDOM() {
        const atChar = this.view.coordsAtPos(this.match.from+1);

        if(!atChar) {
            return;
        }

        let wrapperRect = document.getElementById(EDITOR_WRAPPER_ID).getBoundingClientRect();
        wrapperRect.x += window.scrollX;
        wrapperRect.y += window.scrollY;
        let x = atChar.left;
        let y = atChar.top;
        let height = atChar.bottom-atChar.top;
        x -= wrapperRect.x;
        y -= wrapperRect.y;
        this.dom.style.left = x - 8 + "px";
        this.dom.style.top = y + window.scrollY + height + 8 + "px";
        this.dom.innerHTML = SEARCH_PLACEHOLDER;
        this.dom.style.display = "";
    }

    updateResultsDOM(results) {
        this.dom.innerHTML = "";
        this.results = [];
        let index = 0;

        if(results) {
            if(results.user) {
                for (let user of results.user) {
                    this.appendResultDOM("user", user.organization_member.username, `${user.firstname} ${user.lastname}`, index);
                    index++;
                }
            }

            if(results.groups) {
                for (let group of results.groups) {
                    this.appendResultDOM("group", group.id, group.name, index);
                    index++;
                }
            }

            if(results.articles) {
                for (let article of results.articles) {
                    this.appendResultDOM("article", article.id, article.latest_article_content.title, index);
                    index++;
                }
            }

            if(results.lists) {
                for (let list of results.lists) {
                    this.appendResultDOM("list", list.id, list.name.name, index);
                    index++;
                }
            }

            if(results.tags) {
                for (let tag of results.tags) {
                    this.appendResultDOM("tag", tag.name, tag.name, index);
                    index++;
                }
            }

            this.updateSelectedResultDOM();
        }

        if(!this.results.length) {
            this.dom.innerHTML = SEARCH_PLACEHOLDER;
        }
    }

    appendResultDOM(type, id, content, index) {
        let result = document.createElement("li");
        result.className = `${type}`;
        result.innerText = content;
        result.setAttribute("object", type);
        result.setAttribute("objectid", id);
        this.dom.appendChild(result);
        this.results.push(result);
        result.addEventListener("click", () => {
            this.insert(index);
        });
    }

    handleKeyDown(e) {
        if(!this.active) {
            return false;
        }

        let enter = e.keyCode === 13;
        let down = e.keyCode === 40;
        let up = e.keyCode === 38;
        let left = e.keyCode === 37;
        let right = e.keyCode === 39;
        let esc = e.keyCode === 27;

        if(this.results.length && (enter || down || up || left || right || esc)) {
            e.preventDefault();

            if(enter || right) {
                this.insert(this.selectedResultIndex);
            }
            else if(up) {
                this.selectedResultIndex--;
                this.setSelectedResult();
            }
            else if(down) {
                this.selectedResultIndex++;
                this.setSelectedResult();
            }
            else if(esc) {
                this.close();
            }

            return true;
        }

        return false;
    }

    setSelectedResult() {
        if(this.selectedResultIndex < 0) {
            this.selectedResultIndex = this.results.length-1;
        }
        else if(this.selectedResultIndex > this.results.length-1) {
            this.selectedResultIndex = 0;
        }

        this.updateSelectedResultDOM();
    }

    updateSelectedResultDOM() {
        let nodes = this.dom.childNodes;

        for(let i = 0; i < nodes.length; i++) {
            if(i === this.selectedResultIndex) {
                nodes[i].classList.add("active");
            }
            else {
                nodes[i].classList.remove("active");
            }
        }

        this.scrollIntoView();
    }

    scrollIntoView() {
        if(!this.results.length) {
            return;
        }

        let area = this.dom;
        let height = parseInt(window.getComputedStyle(area, null).getPropertyValue("height"));
        let paddingTop = parseInt(window.getComputedStyle(area, null).getPropertyValue("padding-top"));
        let active = this.results[this.selectedResultIndex];
        let marginTop = parseInt(window.getComputedStyle(active, null).getPropertyValue("margin-top"));
        let start = area.scrollTop;
        let end = start+height;

        if(active.offsetTop <= start || active.offsetTop >= end) {
            area.scrollTop = active.offsetTop - paddingTop - marginTop;
        }
    }

    insert(index) {
        let result = this.results[index];
        let id = result.getAttribute("objectid");
        let type = result.getAttribute("object");
        let title = result.innerText;
        let node = editorSchema.nodes.mention.create({id, type, title, time: new Date().toISOString()});
        this.view.dispatch(this.view.state.tr.replaceWith(this.match.from, this.match.to, node));
        this.view.focus();
        this.close();
    }

    getDOM() {
        return this.dom;
    }

    getDecorator() {
        if(!this.active) {
            return null;
        }

        return DecorationSet.create(this.view.state.doc, [
            Decoration.inline(this.match.from, this.match.to, {
                nodeName: "span",
                class: "mention"
            })
        ]);
    }
}

export function mentionsPlugin() {
    let mentions = new Mentions();
    let plugin = new Plugin({
        view(view) {
            if(!this.mentions) {
                mentions.init(view);
                this.mentions = mentions;
                view.dom.parentNode.insertBefore(this.mentions.getDOM(), view.dom);
            }

            return this.mentions;
        },
        state: {
            init() {
                return mentions;
            },
            apply(tr, mentions) {
                mentions.update(tr);
                return mentions;
            }
        },
        props: {
            handleKeyDown(view, e) {
                let mentions = this.getState(view.state);
                return mentions.handleKeyDown(e);
            },
            decorations(state) {
                let mentions = this.getState(state);
                return mentions.getDecorator();
            }
        }
    });

    plugin.destroy = function() {
        mentions.destroy();
    };

    return plugin;
}
