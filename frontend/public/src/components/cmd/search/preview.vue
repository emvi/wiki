<template>
    <div class="preview">
        <div class="preview-content" ref="preview">
            <emvi-cmd-preview-article v-if="entity.type === 'article'"
                :id="entity.id"
                v-on:loaded="loaded"></emvi-cmd-preview-article>
            <emvi-cmd-preview-article-list v-if="entity.type === 'list'"
                :id="entity.id"
                v-on:loaded="loaded"></emvi-cmd-preview-article-list>
            <emvi-cmd-preview-tag v-if="entity.type === 'tag'"
                :id="entity.id"
                v-on:loaded="loaded"></emvi-cmd-preview-tag>
            <emvi-cmd-preview-group v-if="entity.type === 'group'"
                :id="entity.id"
                v-on:loaded="loaded"></emvi-cmd-preview-group>
            <emvi-cmd-preview-member v-if="entity.type === 'user'"
                :id="entity.id"
                v-on:loaded="loaded"></emvi-cmd-preview-member>
        </div>
    </div>
</template>

<script>
    import {scrollArea} from "../../../util";
    import emviCmdPreviewArticle from "./preview-article.vue";
    import emviCmdPreviewArticleList from "./preview-article-list.vue";
    import emviCmdPreviewTag from "./preview-tag.vue";
    import emviCmdPreviewGroup from "./preview-group.vue";
    import emviCmdPreviewMember from "./preview-member.vue";

    export default {
        components: {
            emviCmdPreviewArticle,
            emviCmdPreviewArticleList,
            emviCmdPreviewTag,
            emviCmdPreviewGroup,
            emviCmdPreviewMember
        },
        props: ["entity", "up", "down"],
        watch: {
            up(up) {
                if(up) {
                    scrollArea(this.$refs.preview, -1);
                }
            },
            down(down) {
                if(down) {
                    scrollArea(this.$refs.preview, 1);
                }
            }
        },
        methods: {
            loaded() {
                this.$emit("loaded");
            }
        }
    }
</script>
