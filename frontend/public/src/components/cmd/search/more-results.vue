<template>
    <div class="entry">
        <emvi-cmd-icon icon="arrow"></emvi-cmd-icon>
        <div class="item">
            <p>{{results}} {{$t("text")}}</p>
        </div>
    </div>
</template>

<script>
    import emviCmdIcon from "../content/icon.vue";

    export default {
        components: {emviCmdIcon},
        props: ["results"]
    }
</script>

<i18n>
    {
        "en": {
            "text": "more results"
        },
        "de": {
            "text": "weitere Ergebnisse"
        }
    }
</i18n>
