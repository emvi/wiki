<template>
    <emvi-label :title="modTime | moment('LLL')">
        {{$t("last_modified")}} {{modTime | moment("from", "now")}} {{$t("modified")}}
    </emvi-label>
</template>

<script>
    import emviLabel from "./label.vue";

    export default {
        components: {emviLabel},
        props: {
            modTime: {default: () => new Date()}
        }
    }
</script>

<i18n>
    {
        "en": {
            "last_modified": "edited",
            "modified": " "
        },
        "de": {
            "last_modified": " ",
            "modified": "bearbeitet"
        }
    }
</i18n>
