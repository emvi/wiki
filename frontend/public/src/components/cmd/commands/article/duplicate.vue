<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-button icon="copy"
                         color="green"
                         :label="$t('label_action')"
                         v-on:enter="run"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {ArticleService} from "../../../../service";
    import {slugWithId} from "../../../../util";
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
            run() {
                this.resetError();
                let id = this.$store.state.page.meta.get("id");
                let langId = this.$store.state.page.meta.get("langId");
                let title = this.$store.state.page.meta.get("title");

                ArticleService.copyArticle(id, langId)
                    .then(newId => {
                        this.$store.dispatch("success", this.$t("toast_duplicated"));
                        this.$store.dispatch("resetCmd");
                        this.$router.push(`/read/${slugWithId(title, newId)}`);
                    })
                    .catch(e => {
                        this.setError(e);
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
            "text": "The article will be duplicated including all translations, the history, tags and notifications. You will be redirected to the copy of this article afterwards.",
            "label_action": "Duplicate",
            "toast_duplicated": "You are now viewing the duplicated article."
        },
        "de": {
            "text": "Der Artikel wird einschließlich aller Übersetzungen, dem Verlauf, Tags und Berechtigungen dupliziert. Du wirst danach zur Kopie weitergeleitet.",
            "label_action": "Duplizieren",
            "toast_duplicated": "Du betrachtest nun den duplizierten Artikel."
        }
    }
</i18n>
