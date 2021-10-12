<template>
    <div class="tags no-select">
        <emvi-tag v-for="tag in tags"
                  :key="tag"
                  :tag="tag"
                  :article-id="articleId"
                  :read-only="readOnly"
                  v-on:remove="remove"></emvi-tag>
        <emvi-add-tag :article-id="articleId"
                      :tags="tags"
                      v-on:add="add"></emvi-add-tag>
    </div>
</template>

<script>
    import emviTag from "./tag.vue";
    import emviAddTag from "./add-tag.vue";

    export default {
        components: {emviTag, emviAddTag},
        props: ["tags", "articleId", "readOnly"],
        methods: {
            add(tag) {
                this.$emit("add", tag);
            },
            remove(tag) {
                this.$emit("remove", tag);
            }
        }
    }
</script>
