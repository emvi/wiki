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
                    {path: "/watch/articles", label: this.$t("tab_articles")},
                    {path: "/watch/lists", label: this.$t("tab_lists")},
                    {path: "/watch/groups", label: this.$t("tab_groups")}
                ]
            };
        }
    }
</script>

<i18n>
    {
        "en": {
            "tab_articles": "Articles",
            "tab_lists": "Lists",
            "tab_groups": "Groups"
        },
        "de": {
            "tab_articles": "Artikel",
            "tab_lists": "Listen",
            "tab_groups": "Gruppen"
        }
    }
</i18n>
