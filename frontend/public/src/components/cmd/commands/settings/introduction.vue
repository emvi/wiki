<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-button icon="help"
                         :label="$t('label_reset')"
                         v-on:enter="run"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {UserService} from "../../../../service";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdButton},
        props: ["esc"],
        data() {
            return {
                name: ""
            };
        },
        watch: {
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        methods: {
            run() {
                UserService.setIntroduction(true)
                    .then(() => {
                        this.$store.dispatch("resetCmd");
                        this.$router.push("/intro");
                    })
                    .catch(e => {
                        this.showTechnicalError(e);
                    });
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
            "text": "You can repeat the introduction. This won't have any affect on this organization. You will be redirected to the introduction on confirmation.",
            "label_reset": "Repeat Introduction"
        },
        "de": {
            "text": "Du kannst die Einführung wiederholen. Diese Aktion hat keinen Einfluss auf die Organisation. Du wirst bei Bestätigung zur Einleitung umgeleitet.",
            "label_reset": "Einführung wiederholen"
        }
    }
</i18n>
