<template>
    <emvi-layout v-on:up="up" v-on:down="down" hide-navigation="true">
        <emvi-enlarge v-show="enlarge" v-on:close="enlarge = false">
            <div ref="enlargeContent"></div>
        </emvi-enlarge>
        <emvi-layout-three-columns>
            <template slot="left">
                <emvi-toc :toc="toc"
                          :list-id="listId"
                          :article-id="id"
                          v-on:jump="jumpToc"></emvi-toc>
            </template>
            <template>
                <div :class="{rtl}" v-if="id">
                    <div class="item-headline">
                        <span class="labels">
                            <emvi-public-status v-if="readEveryone && writeEveryone"></emvi-public-status>
                            <emvi-public-read-status v-if="readEveryone && !writeEveryone"></emvi-public-read-status>
                            <emvi-limited-status v-if="!readEveryone && !writeEveryone && !private"></emvi-limited-status>
                            <emvi-private-status v-if="private"></emvi-private-status>
                            <emvi-authors-label :authors="authors"></emvi-authors-label>
                            <emvi-external-label v-if="clientAccess"></emvi-external-label>
                            <emvi-published-label :published="published"></emvi-published-label>
                            <emvi-modified-label :mod-time="modTime"></emvi-modified-label>
                            <emvi-reading-time-label :reading-time="readingTime"></emvi-reading-time-label>
                            <small class="label uppercase" :title="language.name" v-if="language.code">{{language.code}}</small>
                        </span>
                        <emvi-bookmarked-label type="article" :bookmarked="bookmarked"></emvi-bookmarked-label>
                        <emvi-watched-label type="article" :observed="observed"></emvi-watched-label>
                        <emvi-pinned-label type="article" :pinned="pinned"></emvi-pinned-label>
                    </div>
                    <h1>{{title}}</h1>
                    <emvi-tags :tags="tags"
                               :article-id="id"
                               v-on:add="addTag"
                               v-on:remove="removeTag"></emvi-tags>
                    <emvi-article-archived :reason="archived" v-if="archived"></emvi-article-archived>
                    <emvi-article-version :content="contentObj" v-if="version"></emvi-article-version>
                    <emvi-article-translation-missing v-if="translationMissing"></emvi-article-translation-missing>
                    <div class="article-content"
                         v-html="content"
                         v-on:click="contentClick"
                         ref="content"></div>
                </div>
                <emvi-not-found v-if="!id"></emvi-not-found>
            </template>
        </emvi-layout-three-columns>
    </emvi-layout>
</template>

<script>
    import {mapGetters} from "vuex";
    import {getIdFromSlug, scroll} from "../util";
    import {ArticleService, LangService} from "../service";
    import {addTableEnlargeHandles} from "../editor";
    import {loadCodeMirror} from "../editor/codemirror";
    import {TocMixin} from "./toc";
    import {EnlargeMixin} from "./enlarge";
    import {MentionMixin} from "./mention";
    import {setPageTitle} from "./title";
    import {
        emviLayout,
        emviLayoutThreeColumns,
        emviNotFound,
        emviTags,
        emviToc,
        emviArticleArchived,
        emviArticleVersion,
        emviArticleTranslationMissing,
        emviEnlarge,
        emviAuthorsLabel,
        emviReadingTimeLabel,
        emviViewsLabel,
        emviPublishedLabel,
        emviUnpublishedLabel,
        emviModifiedLabel,
        emviExternalLabel,
        emviPinnedLabel,
        emviBookmarkedLabel,
        emviWatchedLabel,
        emviPublicStatus,
        emviPublicReadStatus,
        emviLimitedStatus,
        emviPrivateStatus
    } from "../components";

    export default {
        mixins: [TocMixin, EnlargeMixin, MentionMixin],
        components: {
            emviLayout,
            emviLayoutThreeColumns,
            emviNotFound,
            emviTags,
            emviToc,
            emviArticleArchived,
            emviArticleVersion,
            emviArticleTranslationMissing,
            emviEnlarge,
            emviAuthorsLabel,
            emviReadingTimeLabel,
            emviViewsLabel,
            emviPublishedLabel,
            emviUnpublishedLabel,
            emviModifiedLabel,
            emviExternalLabel,
            emviPinnedLabel,
            emviBookmarkedLabel,
            emviWatchedLabel,
            emviPublicStatus,
            emviPublicReadStatus,
            emviLimitedStatus,
            emviPrivateStatus
        },
        data() {
            return {
                enlarge: false,
                id: "",
                version: 0,
                title: "",
                contentObj: {},
                content: "",
                rtl: false,
                archived: "",
                observed: false,
                bookmarked: false,
                pinned: false,
                translationMissing: false,
                language: {},
                tags: [],
                authors: [],
                views: 0,
                published: new Date(),
                modTime: new Date(),
                readingTime: null,
                wip: false,
                readEveryone: false,
                writeEveryone: false,
                private: false,
                clientAccess: false,
                listId: ""
            };
        },
        computed: {
            ...mapGetters(["metaUpdate"])
        },
        watch: {
            metaUpdate() {
                let archived = this.$store.state.page.meta.get("archived");
                let pinned = this.$store.state.page.meta.get("pinned");
                let bookmarked = this.$store.state.page.meta.get("bookmarked");
                let observed = this.$store.state.page.meta.get("observed");

                if(archived !== undefined) {
                    this.archived = archived;
                }

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
        mounted() {
            this.id = getIdFromSlug(this.$route.params.slug);
            this.version = this.$route.query.version || 0;
            this.listId = this.$route.query.list || "";
            this.loadArticle(this.$route.query.lang || "");
        },
        methods: {
            loadArticle(langId) {
                ArticleService.getArticle(this.id, langId, this.version)
                    .then(data => {
                        this.$store.dispatch("setMetaVars", [
                            {key: "id", value: this.id},
                            {key: "langId", value: data.content.language_id},
                            {key: "archived", value: data.article.archived},
                            {key: "title", value: data.content.title},
                            {key: "version", value: data.content.version},
                            {key: "observed", value: data.observed},
                            {key: "bookmarked", value: data.bookmarked},
                            {key: "pinned", value: data.article.pinned},
                            {key: "write", value: data.write},
                            {key: "private", value: data.article.private},
                            {key: "translationMissing", value: !data.content.id}
                        ]);

                        this.loadLang(data.content.language_id);
                        this.title = data.content.title || this.$t("translation_missing_title");
                        this.tags = data.article.tags;
                        this.authors = data.authors;
                        this.views = data.article.views;
                        this.published = data.article.published;
                        this.modTime = data.article.mod_time;
                        this.readingTime = data.content.reading_time;
                        this.wip = data.article.wip;
                        this.readEveryone = data.article.read_everyone;
                        this.writeEveryone = data.article.write_everyone;
                        this.private = data.article.private;
                        this.clientAccess = data.article.client_access;
                        this.contentObj = data.content;
                        this.content = addTableEnlargeHandles(data.content.content);
                        this.rtl = data.content.rtl;
                        this.archived = data.article.archived;
                        this.observed = data.observed;
                        this.bookmarked = data.bookmarked;
                        this.pinned = data.article.pinned;
                        this.translationMissing = !data.content.id;
                        setPageTitle(this.title);
                        loadCodeMirror(this.$refs.content);
                        this.buildTableOfContents(this.content);
                        this.showConfirmRecommendations(data.recommendations || []);
                    })
                    .catch(e => {
                        console.error(e);
                        this.id = "";
                        this.$store.dispatch("setMeta", {key: "notFound", value: true});
                    });
            },
            loadLang(lang) {
                LangService.getLang(lang)
                    .then(lang => {
                        this.language = lang;
                        this.$store.dispatch("setMeta", {key: "lang_id", value: this.language.id});
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            addTag(tag) {
                this.tags.push(tag);
            },
            removeTag(tag) {
                let index = this.tags.indexOf(tag);

                if(index > -1) {
                    this.tags.splice(index, 1);
                }
            },
            showConfirmRecommendations(recommendations) {
                if(recommendations.length) {
                    this.$store.dispatch("setMeta", {key: "recommendations", value: recommendations});
                    this.$store.dispatch("resetCmd", this.$t("recommendations_cmd"));
                    this.$store.dispatch("pushColumn", "recommendations");
                }
            },
            contentClick(e) {
                let target = e.target;
                let mention = target.getAttribute("mention");
                let button = e.button;

                if(mention) {
                    e.preventDefault();
                    e.stopPropagation();
                    this.followMention(target, mention, button);
                }
                else if(button === 0 && target.tagName.toLowerCase() === "img" && !target.parentNode.classList.contains("image")) {
                    this.enlargeImage(target);
                }
                else if(button === 0 && target.tagName.toLowerCase() === "div" && target.classList.contains("table-handle")) {
                    this.enlargeTable(target.parentNode);
                }
            },
            up() {
                scroll(-1);
            },
            down() {
                scroll(1);
            }
        }
    }
</script>

<style lang="scss">
    /* hide Codemirror cursor in code blocks */
    .article-content {
        .CodeMirror-cursor {
            display: none;
        }
    }
</style>

<i18n>
    {
        "en": {
            "translation_missing_title": "Translation not found",
            "recommendations_cmd": "/recommendations"
        },
        "de": {
            "translation_missing_title": "Übersetzung nicht gefunden",
            "recommendations_cmd": "/empfehlungen"
        }
    }
</i18n>
