<template>
    <emvi-label :title="$t('title')">{{$t("private")}}</emvi-label>
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
            "private": "Private",
            "title": "Only you can read and edit this article."
        },
        "de": {
            "private": "Privat",
            "title": "Nur du kannst diesen Artikel lesen und bearbeiten."
        }
    }
</i18n>
