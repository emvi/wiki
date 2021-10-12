<template>
    <emvi-label>
        <span v-if="views > 1">{{views}} {{$t("views")}}</span>
        <span v-if="views === 1">{{views}} {{$t("view")}}</span>
    </emvi-label>
</template>

<script>
    import emviLabel from "./label.vue";

    export default {
        components: {emviLabel},
        props: ["views"]
    }
</script>

<i18n>
    {
        "en": {
            "views": "views",
            "view": "view"
        },
        "de": {
            "views": "Aufrufe",
            "view": "Aufruf"
        }
    }
</i18n>
