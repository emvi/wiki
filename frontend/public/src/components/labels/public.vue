<template>
    <emvi-label :title="$t('title')">{{$t("public")}}</emvi-label>
</template>

<script>
    import emviLabel from "./label.vue";

    export default {
        components: {emviLabel}
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
