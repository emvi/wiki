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
    import {UsergroupService} from "../../../../service";
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

                UsergroupService.deleteUsergroup(id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_deleted"));
                        this.$store.dispatch("resetCmd");
                        this.$router.push("/groups");
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
            "text": "Are you sure you want to delete this group?",
            "label_delete": "Delete Group",
            "toast_deleted": "The group has been deleted."
        },
        "de": {
            "text": "Bist du sicher, dass du diese Gruppe löschen möchtest?",
            "label_delete": "Gruppe löschen",
            "toast_deleted": "Die Gruppe wurde gelöscht."
        }
    }
</i18n>
