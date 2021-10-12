<template>
    <emvi-status icon="lock" :title="$t('title')">
        {{$t("private")}}
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
            "private": "Private",
            "title": "Only you can read and edit this article."
        },
        "de": {
            "private": "Privat",
            "title": "Nur du kannst diesen Artikel lesen und bearbeiten."
        }
    }
</i18n>
