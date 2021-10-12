<template>
    <emvi-preview-card icon="tag" :active="active" :up="up" :down="down" :preview="showPreview" :scroll-area="scrollArea" v-on:click="open" v-on:preview="loadTag" ref="card">
        <template>
            <span v-html="title"></span>
        </template>
        <template slot="info">
            <span :title="entity.def_time | moment('LT')">{{entity.def_time | moment("ll")}}</span>
        </template>
        <template slot="preview" v-if="showPreview">
            <emvi-preview-tag :id="entity.id" v-on:loaded="scroll"></emvi-preview-tag>
        </template>
    </emvi-preview-card>
</template>

<script>
    import {markInText} from "../cmd/util";
    import {CardMixin} from "./preview-card";
    import emviPreviewCard from "./preview-card.vue";
    import emviPreviewTag from "./preview-tag.vue";

    export default {
        mixins: [CardMixin],
        components: {emviPreviewCard, emviPreviewTag},
        computed: {
            title() {
                return markInText(this.query, this.entity.name)+` (${this.entity.usages})`;
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
                    this.loadTag();
                }
            }
        },
        methods: {
            loadTag() {
                this.showPreview = !this.showPreview;
                this.$emit("preview", {show: this.showPreview, index: this.index});
            },
            open() {
                this.$store.dispatch("select", 0);
                this.$router.push(`/tag/${this.entity.name}`);
            }
        }
    }
</script>
