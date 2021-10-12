<template>
    <emvi-tabs :tabs="tabs"></emvi-tabs>
</template>

<script>
    import emviTabs from "./tabs.vue";

    export default {
        components: {emviTabs},
        data() {
            return {
                tabs: [
                    {path: "/bookmarks/articles", label: this.$t("tab_articles")},
                    {path: "/bookmarks/lists", label: this.$t("tab_lists")}
                ]
            };
        }
    }
</script>

<i18n>
    {
        "en": {
            "tab_articles": "Articles",
            "tab_lists": "Lists"
        },
        "de": {
            "tab_articles": "Artikel",
            "tab_lists": "Listen"
        }
    }
</i18n>
