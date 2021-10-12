<template>
    <emvi-preview-card icon="article" :active="active" :disabled="entity.archived" :scroll-area="scrollArea" :up="up" :down="down" :preview="showPreview" v-on:click="open" v-on:preview="loadArticle" ref="card">
        <template>
            <span v-html="title"></span>
            <emvi-private-label v-if="entity.private"></emvi-private-label>
            <emvi-external-label v-if="entity.client_access"></emvi-external-label>
        </template>
        <template slot="info">
            <span :title="entity | authorNamesFilter">
                {{entity.latest_article_content.authors[0].firstname}} {{entity.latest_article_content.authors[0].lastname}}
                <template v-if="entity.latest_article_content.authors.length > 1">+ {{entity.latest_article_content.authors.length-1}}</template>
            </span>
            <span class="dot">·</span>
            <span>{{entity.views}} {{$t("views")}}</span>
            <span class="dot">·</span>
            <span :title="entity.published | moment('LT')" v-if="entity.published">{{entity.published | moment("ll")}}</span>
            <span v-if="!entity.published">{{$t("unpublished")}}</span>
            <span class="dot">·</span>
            <span :title="entity.latest_article_content.mod_time | moment('LLL')">{{$t("edited_before")}} {{entity.latest_article_content.mod_time | moment("from", "now")}} {{$t("edited_after")}}</span>
        </template>
        <template slot="preview" v-if="showPreview">
            <emvi-preview-article :id="entity.id" v-on:loaded="scroll"></emvi-preview-article>
        </template>
    </emvi-preview-card>
</template>

<script>
    import {authorNamesFilter, slugWithId} from "../../util";
    import {markInText} from "../cmd/util";
    import {CardMixin} from "./preview-card";
    import emviPreviewCard from "./preview-card.vue";
    import emviPreviewArticle from "./preview-article.vue";
    import emviPrivateLabel from "../labels/private.vue";
    import emviExternalLabel from "../labels/external.vue";

    export default {
        mixins: [CardMixin],
        components: {emviPreviewCard, emviPreviewArticle, emviPrivateLabel, emviExternalLabel},
        filters: {
            authorNamesFilter
        },
        computed: {
            title() {
                return markInText(this.query, this.entity.latest_article_content.title);
            }
        },
        watch: {
            enter(enter) {
                if(enter && this.active) {
                    this.open();
                }
            },
            tab(tab) {
                if(tab && this.active) {
                    this.loadArticle();
                }
            }
        },
        methods: {
            loadArticle() {
                this.showPreview = !this.showPreview;
                this.$emit("preview", {show: this.showPreview, index: this.index});
            },
            open() {
                this.$store.dispatch("select", 0);
                let listId = this.entity.list_id ? `?list=${this.entity.list_id}` : "";
                this.$router.push(`/read/${slugWithId(this.entity.latest_article_content.title, this.entity.id)}${listId}`);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "views": "views",
            "unpublished": "Unpublished",
            "edited_before": "edited",
            "edited_after": " "
        },
        "de": {
            "views": "Aufrufe",
            "unpublished": "Unveröffentlicht",
            "edited_before": " ",
            "edited_after": "bearbeitet"
        }
    }
</i18n>
