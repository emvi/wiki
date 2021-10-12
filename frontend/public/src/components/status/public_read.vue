<template>
    <emvi-status icon="view" :title="$t('title')">
        {{$t("public_read")}}
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
            "public_read": "Public",
            "title": "Anyone can read this article. Only selected members can edit it."
        },
        "de": {
            "public_read": "Öffentlich",
            "title": "Jeder kann diesen Artikel lesen. Nur ausgewählte Mitglieder können ihn bearbeiten."
        }
    }
</i18n>
