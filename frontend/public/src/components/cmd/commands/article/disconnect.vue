<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-button icon="refresh"
                         color="green"
                         :label="$t('label_action')"
                         v-on:enter="reload"
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
            reload() {
                location.reload();
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
            "text": "The connection to the collaboration server was lost. Your changes have been saved. Please reload the article to reconnect. In case you have created a new article you can find the draft in your library. For an existing article you can click back in your browser.",
            "label_action": "Reload Page"
        },
        "de": {
            "text": "Die Verbindung zum Kollaborationsserver wurde unterbrochen. Deine Änderungen wurden gespeichert. Bitte lade den Artikel neu, um dich erneut zu verbinden. Wenn du einen neuen Artikel angelegt hast findest du ihn in deiner Bibliothek, für existierende Artikel kannst du einfach auf zurück in deinem Browser klicken.",
            "label_action": "Seite Neuladen"
        }
    }
</i18n>
