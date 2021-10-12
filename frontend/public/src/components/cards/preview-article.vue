<template>
    <div>
        <emvi-loading v-show="!loaded"></emvi-loading>
        <emvi-preview-empty icon="article" v-show="loaded && !content"></emvi-preview-empty>
        <div class="article-content" v-html="content" v-show="content"></div>
    </div>
</template>

<script>
    import {ArticleService} from "../../service";
    import emviPreviewEmpty from "./preview-empty.vue";
    import emviLoading from "../content/loading.vue";

    export default {
        components: {emviPreviewEmpty, emviLoading},
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
                    this.resetError();
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
