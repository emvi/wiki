<template>
    <div>
        <div v-if="!listId">
            <emvi-toc-headline v-for="headline in toc"
                :headline="headline"
                :key="headline.id"
                v-on:jump="jump"></emvi-toc-headline>
        </div>
        <div v-else>
            <emvi-read-list :list-id="listId" :article-id="articleId">
                <emvi-toc-headline v-for="headline in toc"
                    :headline="headline"
                    :key="headline.id"
                    v-on:jump="jump"></emvi-toc-headline>
            </emvi-read-list>
        </div>
    </div>
</template>

<script>
    import emviTocHeadline from "./toc-headline.vue";
    import emviReadList from "./read-list.vue";

    export default {
        components: {emviTocHeadline, emviReadList},
        props: ["toc", "listId", "articleId"],
        methods: {
            jump(headline) {
                this.$emit("jump", headline);
            }
        }
    }
</script>
