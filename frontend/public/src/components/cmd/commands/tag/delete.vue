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
    import {TagService} from "../../../../service";
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

                TagService.deleteTag(id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_deleted"));
                        this.$store.dispatch("resetCmd");
                        this.$router.push("/tags");
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
            "text": "Are you sure you want to delete this tag?",
            "label_delete": "Delete Tag",
            "toast_deleted": "The tag has been deleted."
        },
        "de": {
            "text": "Bist du sicher, dass du diesen Tag löschen möchtest?",
            "label_delete": "Tag löschen",
            "toast_deleted": "Der Tag wurde gelöscht."
        }
    }
</i18n>
