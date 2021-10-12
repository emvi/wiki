<template>
    <emvi-list :title="list.name.name"
               title-icon="list"
               icon="article"
               card="emvi-article-card"
               :perform-search="search"
               disable-search="true"
               v-on:open="open">
        <template slot="headline">
            <span class="labels">
                <small v-if="list.name.info">{{list.name.info}}</small>
                <emvi-external-label v-if="list.client_access"></emvi-external-label>
                <emvi-published-label :published="list.def_time"></emvi-published-label>
                <emvi-modified-label :mod-time="list.mod_time"></emvi-modified-label>
            </span>
            <emvi-bookmarked-label type="list" :bookmarked="bookmarked"></emvi-bookmarked-label>
            <emvi-watched-label type="list" :observed="observed"></emvi-watched-label>
            <emvi-pinned-label type="list" :pinned="pinned"></emvi-pinned-label>
        </template>
    </emvi-list>
</template>

<script>
    import {mapGetters} from "vuex";
    import {ArticlelistService} from "../service";
    import {addAttrToListElements, getIdFromSlug, slugWithId} from "../util";
    import {setPageTitle} from "./title";
    import {
        emviList,
        emviExternalLabel,
        emviPublishedLabel,
        emviModifiedLabel,
        emviPinnedLabel,
        emviBookmarkedLabel,
        emviWatchedLabel
    } from "../components";

    export default {
        components: {
            emviList,
            emviExternalLabel,
            emviPublishedLabel,
            emviModifiedLabel,
            emviPinnedLabel,
            emviBookmarkedLabel,
            emviWatchedLabel
        },
        data() {
            return {
                list: {
                    name: {},
                    def_time: new Date(),
                    mod_time: new Date()
                },
                observed: false,
                bookmarked: false,
                pinned: false
            };
        },
        computed: {
            ...mapGetters(["metaUpdate"])
        },
        watch: {
            metaUpdate(update) {
                if(update > 1 && this.$store.state.page.meta.get("updated")) {
                    this.list.name.name = this.$store.state.page.meta.get("name");
                    this.list.name.info = this.$store.state.page.meta.get("info");
                    this.list.public = this.$store.state.page.meta.get("public");
                    this.list.client_access = this.$store.state.page.meta.get("client_access");
                    this.$store.dispatch("setMeta", {key: "updated", value: undefined});
                }

                let pinned = this.$store.state.page.meta.get("pinned");
                let bookmarked = this.$store.state.page.meta.get("bookmarked");
                let observed = this.$store.state.page.meta.get("observed");

                if(pinned !== undefined) {
                    this.pinned = pinned;
                }

                if(bookmarked !== undefined) {
                    this.bookmarked = bookmarked;
                }

                if(observed !== undefined) {
                    this.observed = observed;
                }
            }
        },
        methods: {
            open(result) {
                let listId = this.list.id ? `?list=${this.list.id}` : "";
                this.$router.push(`/read/${slugWithId(result.latest_article_content.title, result.id)}${listId}`);
            },
            search(query, filter, cancelToken) {
                if(!this.list.id) {
                    return this.loadList()
                        .then(() => {
                            return this.loadArticles(query, filter, cancelToken);
                        })
                        .catch(e => {
                            return Promise.reject(e);
                        });
                }
                else {
                    return this.loadArticles(query, filter, cancelToken);
                }
            },
            loadList() {
                let id = getIdFromSlug(this.$route.params.slug);

                return new Promise((resolve, reject) => {
                    ArticlelistService.getArticlelist(id)
                        .then(list => {
                            this.$store.dispatch("setMetaVars", [
                                {key: "id", value: list.list.id},
                                {key: "observed", value: list.observed},
                                {key: "bookmarked", value: list.bookmarked},
                                {key: "pinned", value: list.pinned},
                                {key: "moderator", value: list.moderator}
                            ]);
                            this.observed = list.observed;
                            this.bookmarked = list.bookmarked;
                            this.pinned = list.list.pinned;
                            this.list = list.list;
                            setPageTitle(this.list.name.name);
                            resolve();
                        })
                        .catch(e => {
                            e.status = 404;
                            reject(e);
                        });
                });
            },
            loadArticles(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    ArticlelistService.getEntries(this.list.id, filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results: addAttrToListElements(results, "list_id", this.list.id), count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "edited_before": "edited",
            "edited_after": " "
        },
        "de": {
            "edited_before": " ",
            "edited_after": "bearbeitet"
        }
    }
</i18n>
