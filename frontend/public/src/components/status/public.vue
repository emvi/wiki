<template>
    <emvi-status icon="edit" :title="$t('title')">
        {{$t("public")}}
    </emvi-status>
</template>

<script>
    import emviStatus from "./status.vue";

    export default {
        components: {emviStatus}
    }
</script>

<i18n>
    {
        "en": {
            "public": "Public",
            "title": "The article can be read and edited by anyone."
        },
        "de": {
            "public": "Öffentlich",
            "title": "Der Artikel kann von jedem lesen und bearbeitet werden."
        }
    }
</i18n>
