<template>
    <emvi-label>{{$t("unpublished")}}</emvi-label>
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
            "unpublished": "Has unpublished changes"
        },
        "de": {
            "unpublished": "Hat unveröffentlicht Änderungen"
        }
    }
</i18n>
