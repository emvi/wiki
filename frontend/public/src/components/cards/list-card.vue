<template>
    <emvi-preview-card icon="list" :active="active" :up="up" :down="down" :preview="showPreview" :scroll-area="scrollArea" v-on:click="open" v-on:preview="loadList" ref="card">
        <template>
            <span v-html="title"></span>
            <emvi-external-label v-if="entity.client_access"></emvi-external-label>
        </template>
        <template slot="info">
            <span>{{entity.name.info}}</span>
            <span class="dot" v-if="entity.name.info">·</span>
            <span :title="entity.def_time | moment('LT')">{{entity.def_time | moment("ll")}}</span>
            <span class="dot">·</span>
            <span :title="entity.mod_time | moment('LLL')">{{$t("edited_before")}} {{entity.mod_time | moment("from", "now")}} {{$t("edited_after")}}</span>
        </template>
        <template slot="preview" v-if="showPreview">
            <emvi-preview-article-list :id="entity.id" v-on:loaded="scroll"></emvi-preview-article-list>
        </template>
    </emvi-preview-card>
</template>

<script>
    import {slugWithId} from "../../util";
    import {markInText} from "../cmd/util";
    import {CardMixin} from "./preview-card";
    import emviPreviewCard from "./preview-card.vue";
    import emviPreviewArticleList from "./preview-article-list.vue";
    import emviExternalLabel from "../labels/external.vue";

    export default {
        mixins: [CardMixin],
        components: {emviPreviewCard, emviPreviewArticleList, emviExternalLabel},
        computed: {
            title() {
                return markInText(this.query, this.entity.name.name)+` (${this.entity.articles})`;
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
                    this.loadList();
                }
            }
        },
        methods: {
            loadList() {
                this.showPreview = !this.showPreview;
                this.$emit("preview", {show: this.showPreview, index: this.index});
            },
            open() {
                this.$store.dispatch("select", 0);
                this.$router.push(`/list/${slugWithId(this.entity.name.name, this.entity.id)}`);
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
