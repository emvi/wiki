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
    import {ArticlelistService} from "../../../../service";
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

                ArticlelistService.deleteArticlelist(id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_deleted"));
                        this.$store.dispatch("resetCmd");
                        this.$router.push("/lists");
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
            "text": "Are you sure you want to delete this list?",
            "label_delete": "Delete List",
            "toast_deleted": "The list has been deleted."
        },
        "de": {
            "text": "Bist du sicher, dass du diese Liste löschen möchtest?",
            "label_delete": "Liste löschen",
            "toast_deleted": "Die Liste wurde gelöscht."
        }
    }
</i18n>
