<template>
    <div class="toc">
        <router-link class="list" :to="link">
            <i class="icon icon-list green-100 bg-green-10 size-32"></i>
            <h3>{{list.name.name}} ({{count}})</h3>
        </router-link>
        <!-- this is slightly different from the shortcut component, as it does not allow clicking on it -->
        <div class="shortcut no-select" v-if="!isTouch">
            <span>
                <span class="key">{{$t("shift")}}</span>
                <span>&#43;</span>
                <span class="key">&uarr;</span>
                <span class="key">&darr;</span>
                <span>{{$t("shortcut_switch")}}</span>
            </span>
        </div>
        <emvi-read-list-entry v-for="(entry, index) in entries"
                          :key="entry.id"
                          :entry="entry"
                          :pos="pos+index"
                          :list-id="listId"
                          :article-id="articleId">
            <slot></slot>
        </emvi-read-list-entry>
    </div>
</template>

<script>
    import {findIndexById, slugWithId} from "../../util";
    import {ArticlelistService} from "../../service";
    import emviReadListEntry from "./read-list-entry.vue";

    const centerBefore = 4;

    export default {
        components: {emviReadListEntry},
        props: ["listId", "articleId"],
        data() {
            return {
                keydownHandler: null,
                list: {name: {name: ""}},
                entries: [],
                count: 0,
                pos: 0
            };
        },
        computed: {
            link() {
                return `/list/${slugWithId(this.list.name.name, this.list.id)}`;
            }
        },
        mounted() {
            this.bindKeys();
            this.loadList();
        },
        beforeDestroy() {
            this.unbindKeys();
        },
        methods: {
            bindKeys() {
                this.keydownHandler = e => {
                    if(!e.shiftKey) {
                        return;
                    }

                    let prevent = true;

                    switch(e.code) {
                        case "ArrowUp":
                            this.previousArticle();
                            break;
                        case "ArrowDown":
                            this.nextArticle();
                            break;
                        default:
                            prevent = false;
                    }

                    if(prevent) {
                        e.preventDefault();
                        e.stopPropagation();
                    }
                };
                window.addEventListener("keydown", this.keydownHandler);
            },
            unbindKeys() {
                window.removeEventListener("keydown", this.keydownHandler);
            },
            previousArticle() {
                let index = findIndexById(this.entries, this.articleId);

                if(index > -1) {
                    index--;

                    if(index < 0) {
                        return;
                    }

                    this.openArticle(index);
                }
            },
            nextArticle() {
                let index = findIndexById(this.entries, this.articleId);

                if(index > -1) {
                    index++;

                    if(index > this.count-1) {
                        return;
                    }

                    this.openArticle(index);
                }
            },
            openArticle(index) {
                let entry = this.entries[index];
                this.$router.push(`/read/${slugWithId(entry.latest_article_content.title, entry.id)}?list=${this.listId}`);
            },
            loadList() {
                ArticlelistService.getArticlelist(this.listId)
                    .then(data => {
                        this.list = data.list;
                        this.loadListEntries();
                    });
            },
            loadListEntries() {
                let filter = {
                    center_article_id: this.articleId,
                    center_before: centerBefore,
                    limit: centerBefore*2+1
                };

                ArticlelistService.getEntries(this.list.id, filter)
                    .then(data => {
                        this.entries = data.results;
                        this.count = data.count;
                        this.pos = data.start_pos;
                    });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "shortcut_switch": "Switch Article",
            "shift": "Shift"
        },
        "de": {
            "shortcut_switch": "Artikel wechseln",
            "shift": "Shift"
        }
    }
</i18n>
