<template>
    <emvi-preview-card :entity="entity" :active="active" :up="up" :down="down" :preview="showPreview" :scroll-area="scrollArea" v-on:click="open" v-on:preview="loadMember" ref="card">
        <template>
            <span v-html="title"></span>
        </template>
        <template slot="info">
            <span>{{entity.organization_member.username}}</span>
            <span class="dot" v-if="entity.organization_member.info">·</span>
            <span>{{entity.organization_member.info}}</span>
            <span class="dot">·</span>
            <span>{{entity.email}}</span>
        </template>
        <template slot="preview" v-if="showPreview">
            <emvi-preview-member :id="entity.id" v-on:loaded="scroll"></emvi-preview-member>
        </template>
    </emvi-preview-card>
</template>

<script>
    import {markInText} from "../cmd/util";
    import {CardMixin} from "./preview-card";
    import emviPreviewCard from "./preview-card.vue";
    import emviPreviewMember from "./preview-member.vue";

    export default {
        mixins: [CardMixin],
        components: {emviPreviewCard, emviPreviewMember},
        computed: {
            title() {
                return markInText(this.query, `${this.entity.firstname} ${this.entity.lastname}`);
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
                    this.loadMember();
                }
            }
        },
        methods: {
            loadMember() {
                this.showPreview = !this.showPreview;
                this.$emit("preview", {show: this.showPreview, index: this.index});
            },
            open() {
                this.$store.dispatch("select", 0);
                this.$router.push(`/member/${this.entity.organization_member.username}`);
            }
        }
    }
</script>
