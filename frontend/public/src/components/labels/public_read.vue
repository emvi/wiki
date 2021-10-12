<template>
    <emvi-label :title="$t('title')">{{$t("public_read")}}</emvi-label>
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
            "public_read": "Public (read only)",
            "title": "Anyone can read this article. Only selected members can edit it."
        },
        "de": {
            "public_read": "Öffentlich (nur lesen)",
            "title": "Jeder kann diesen Artikel lesen. Nur ausgewählte Mitglieder können ihn bearbeiten."
        }
    }
</i18n>
