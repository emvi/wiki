<template>
    <emvi-cmd-preview-list :perform-search="search"
        icon="tag"
        result-icon="article"
        v-on:loaded="loaded"></emvi-cmd-preview-list>
</template>

<script>
    import {addAttrToListElements} from "../../../util";
    import {SearchService} from "../../../service";
    import emviCmdPreviewList from "./preview-list.vue";

    export default {
        components: {emviCmdPreviewList},
        props: ["id"],
        methods: {
            search() {
                return new Promise((resolve, reject) => {
                    SearchService.findArticles("", {tag_ids: this.id, sort_title: "asc"})
                        .then(({results, count}) => {
                            resolve({results: addAttrToListElements(results, "type", "article"), count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loaded() {
                this.$emit("loaded");
            }
        }
    }
</script>
