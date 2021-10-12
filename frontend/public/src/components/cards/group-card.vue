<template>
    <emvi-preview-card icon="group" :active="active" :up="up" :down="down" :preview="showPreview" :scroll-area="scrollArea" v-on:click="open" v-on:preview="loadGroup" ref="card">
        <template>
            <span v-html="title"></span>
        </template>
        <template slot="info">
            <span>{{entity.info}}</span>
            <span class="dot" v-if="entity.info">·</span>
            <span :title="entity.def_time | moment('LT')">{{entity.def_time | moment("ll")}}</span>
            <span class="dot">·</span>
            <span :title="entity.mod_time | moment('LLL')">{{$t("edited_before")}} {{entity.mod_time | moment("from", "now")}} {{$t("edited_after")}}</span>
        </template>
        <template slot="preview" v-if="showPreview">
            <emvi-preview-group :id="entity.id" v-on:loaded="scroll"></emvi-preview-group>
        </template>
    </emvi-preview-card>
</template>

<script>
    import {slugWithId} from "../../util";
    import {markInText} from "../cmd/util";
    import {CardMixin} from "./preview-card";
    import emviPreviewCard from "./preview-card.vue";
    import emviPreviewGroup from "./preview-group.vue";

    export default {
        mixins: [CardMixin],
        components: {emviPreviewCard, emviPreviewGroup},
        computed: {
            title() {
                return markInText(this.query, this.entity.name)+` (${this.entity.member})`;
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
                    this.loadGroup();
                }
            }
        },
        methods: {
            loadGroup() {
                this.showPreview = !this.showPreview;
                this.$emit("preview", {show: this.showPreview, index: this.index});
            },
            open() {
                this.$store.dispatch("select", 0);
                this.$router.push(`/group/${slugWithId(this.entity.name, this.entity.id)}`);
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
