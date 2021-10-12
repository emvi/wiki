<template>
    <emvi-status icon="unlock" :title="$t('title')">
        {{$t("limited")}}
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
            "limited": "Limited",
            "title": "The article can be read and edited by authorized members only."
        },
        "de": {
            "limited": "Limitiert",
            "title": "Der Artikel kann nur von berechtigten Nutzern gelesen und bearbeitet werden."
        }
    }
</i18n>
