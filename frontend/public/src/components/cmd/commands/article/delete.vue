<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-button icon="trash"
                         color="red"
                         :label="$t('label_delete')"
                         v-on:enter="run"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {ArticleService, UsergroupService} from "../../../../service";
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

                ArticleService.archiveArticle(id, "", true)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_deleted"));
                        this.$store.dispatch("resetCmd");
                        this.$router.push("/articles");
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
            "text": "Are you sure you want to delete this article? The action cannot be undone!",
            "label_delete": "Delete Article",
            "toast_deleted": "The article has been deleted."
        },
        "de": {
            "text": "Bist du sicher, dass du diesen Artikel löschen möchtest? Die Aktion kann nicht rückgängig gemacht werden!",
            "label_delete": "Artikel löschen",
            "toast_deleted": "Der Artikel wurde gelöscht."
        }
    }
</i18n>
