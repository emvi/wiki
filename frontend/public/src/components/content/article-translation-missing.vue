<template>
    <div class="hint red">
        <div class="hint-content">
            <p>{{$t("text")}}</p>
            <small>{{$t("hint")}}</small>
            <template v-if="!isExpert">
                <br />
                <small><strong>{{$t("expert")}}</strong></small>
            </template>
        </div>
        <i v-on:click.stop.prevent="cmdLanguage" class="icon icon-globe size-40 cursor-pointer"></i>
    </div>
</template>

<script>
    export default {
        props: ["reason"],
        methods: {
            cmdLanguage() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("language_cmd"));
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text": "The translation for the selected language was not found.",
            "hint": "Use /language to switch the active language.",
            "language_cmd": "/language",
            "expert": "You need an Expert organization to add translations."
        },
        "de": {
            "text": "Die Übersetzung für die ausgewählte Sprache wurde nicht gefunden.",
            "hint": "Verwende /sprache um die aktive Sprache zu ändern.",
            "language_cmd": "/sprache",
            "expert": "Du benötigst eine Expert Organisation um Übersetzungen hinzuzufügen."
        }
    }
</i18n>
