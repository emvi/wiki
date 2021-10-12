<template>
    <div class="user-list" v-if="authors.length > 0">
        <emvi-author v-for="author in authors"
                     :key="author.id"
                     :author="author"></emvi-author>
    </div>
</template>

<script>
    import emviAuthor from "./author.vue";

    export default {
        components: {emviAuthor},
        props: ["authors"]
    }
</script>
