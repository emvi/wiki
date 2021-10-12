<template>
    <div class="hint blue">
        <div class="hint-content">
            <p>{{$t("text_reason")}} {{content.def_time | moment("LLL")}} {{$t("text_published")}}</p>
            <small>{{$t("label_history")}}</small>
        </div>
        <i v-on:click.stop.prevent="cmd" class="icon icon-history size-40 cursor-pointer"></i>
    </div>
</template>

<script>

    export default {
        props: ["content"],
        methods: {
            cmd() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("history_cmd"));
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text_reason": "This version was published on",
            "text_published": " ",
            "label_history": "Use /history to view and switch versions.",
            "history_cmd": "/history"
        },
        "de": {
            "text_reason": "Diese Version wurde am",
            "text_published": "veröffentlicht",
            "label_history": "Verwende /historie um die Versionen einzusehen und zu wechseln.",
            "history_cmd": "/historie"
        }
    }
</i18n>
