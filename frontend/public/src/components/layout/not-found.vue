<template>
    <div>
        <emvi-page-title disable-result-count="true">{{$t("title")}}</emvi-page-title>
        <p>{{$t("text")}}</p>
    </div>
</template>

<script>
    import emviPageTitle from "../content/page-title.vue";

    export default {
        components: {emviPageTitle}
    }
</script>

<i18n>
    {
        "en": {
            "title": "Not found.",
            "text": "The element has not been found."
        },
        "de": {
            "title": "Nicht gefunden.",
            "text": "Das Element wurde nicht gefunden."
        }
    }
</i18n>
