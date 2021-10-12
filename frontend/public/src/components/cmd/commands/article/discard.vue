<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-button icon="trash"
                         color="red"
                         :label="$t('label_discard')"
                         v-on:enter="discard"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdButton},
        props: ["esc"],
        watch: {
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        methods: {
            discard() {
                let page = this.$store.state.page.meta.get("page");
                page.leave(true);
                this.$store.dispatch("resetCmd");
            },
            cancel() {
                this.$store.dispatch("popColumn");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text": "Are you sure you want to discard the changes?",
            "label_discard": "Discard Changes"
        },
        "de": {
            "text": "Bist du sicher, dass du die Änderungen verwerfen möchtest?",
            "label_discard": "Änderungen verwerfen"
        }
    }
</i18n>
