<template>
    <div>
        <emvi-cmd-loading v-show="!loaded"></emvi-cmd-loading>
        <emvi-cmd-preview-empty icon="article" v-show="loaded && !content"></emvi-cmd-preview-empty>
        <div class="article-content" v-html="content" v-show="content"></div>
    </div>
</template>

<script>
    import {ArticleService} from "../../../service";
    import emviCmdPreviewEmpty from "./preview-empty.vue";
    import emviCmdLoading from "../content/loading.vue";

    export default {
        components: {emviCmdPreviewEmpty, emviCmdLoading},
        props: ["id"],
        data() {
            return {
                content: "",
                loaded: false
            };
        },
        mounted() {
            this.loadArticle();
        },
        methods: {
            loadArticle() {
                if(!this.loaded) {
                    ArticleService.getArticlePreview(this.id)
                        .then(content => {
                            this.content = content;
                            this.loaded = true;
                            this.$emit("loaded");
                        })
                        .catch(e => {
                            this.setError(e);
                        });
                }
            }
        }
    }
</script>
