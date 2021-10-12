<template>
    <div>
        <p v-if="archived">{{$t("text")}}</p>
        <emvi-cmd-input v-if="!archived"
                        :label="$t('label_message')"
                        :index="0"
                        :error="validation['message']"
                        v-model="message"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="run"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-button :icon="archived ? 'back' : 'archive'"
                         :color="archived ? 'green' : 'red'"
                         :label="archived ? $t('label_restore') : $t('label_archive')"
                         :index="1"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="run"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {ArticleService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdInput, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                archived: false,
                message: ""
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 2, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        mounted() {
            this.archived = !!this.$store.state.page.meta.get("archived");

            if(this.archived) {
                this.$store.dispatch("selectRow", 1);
            }
        },
        methods: {
            run() {
                this.resetError();
                let id = this.$store.state.page.meta.get("id");

                ArticleService.archiveArticle(id, this.message, false)
                    .then(() => {
                        this.$store.dispatch("setMeta", {key: "archived", value: this.archived ? "" : this.message});
                        this.$store.dispatch("success", this.$t(this.archived ? "toast_archived_off" : "toast_archived_on"));
                        this.$store.dispatch("resetCmd");
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.$store.dispatch("popColumn");
            },
            nextRow() {
                this.$store.dispatch("selectNextRow");
            },
            previousRow() {
                this.$store.dispatch("selectPreviousRow");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text": "Confirm to restore this article.",
            "label_message": "Reason",
            "label_archive": "Archive",
            "label_restore": "Restore",
            "toast_archived_on": "The article has been archived.",
            "toast_archived_off": "The article has been restored."
        },
        "de": {
            "text": "Bestätige um den Artikel wiederherzustellen.",
            "label_message": "Grund",
            "label_archive": "Archivieren",
            "label_restore": "Wiederherstellen",
            "toast_archived_on": "Der Artikel wurde archiviert.",
            "toast_archived_off": "Der Artikel wurde wiederhergestellt."
        }
    }
</i18n>
