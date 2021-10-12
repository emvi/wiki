<template>
    <div class="entry">
        <emvi-cmd-icon icon="close"></emvi-cmd-icon>
        <div class="item">
            <p>{{$t("text")}}</p>
        </div>
    </div>
</template>

<script>
    import emviCmdIcon from "../content/icon.vue";

    export default {
        components: {emviCmdIcon}
    }
</script>

<i18n>
    {
        "en": {
            "text": "No results found"
        },
        "de": {
            "text": "Keine Ergebnisse gefunden"
        }
    }
</i18n>
