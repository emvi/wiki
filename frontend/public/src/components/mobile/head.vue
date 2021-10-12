<template>
    <div class="mobile-header">
        <emvi-orga></emvi-orga>
        <emvi-avatar v-on:click="openProfile"></emvi-avatar>
    </div>
</template>

<script>
    import emviOrga from "../layout/orga.vue";
    import emviAvatar from "../layout/avatar.vue";

    export default {
        components: {emviOrga, emviAvatar},
        methods: {
            openProfile() {
                window.open(`${EMVI_WIKI_WEBSITE_HOST}/account`, "_blank");
            }
        }
    }
</script>
