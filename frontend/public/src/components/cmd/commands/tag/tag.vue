<template>
    <div>
        <emvi-cmd-input :label="$t('label_name')"
                        :index="0"
                        :error="validation['tag']"
                        v-model="name"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="run"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-button icon="edit"
                         color="green"
                         :label="$t('label_action')"
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
    import {TagService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdInput, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                id: "",
                name: ""
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
            this.id = this.$store.state.page.meta.get("id");
            this.name = this.$store.state.page.meta.get("name");
        },
        methods: {
            run() {
                this.resetError();

                TagService.renameTag(this.id, this.name)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_renamed"));
                        this.$store.dispatch("resetCmd");
                        this.$router.push(`/tag/${this.name}`);
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
            "label_name": "Tag",
            "label_action": "Rename Tag",
            "toast_renamed": "The tag has been renamed."
        },
        "de": {
            "label_name": "Tag",
            "label_action": "Tag umbenennen",
            "toast_renamed": "Der Tag wurde umbenannt."
        }
    }
</i18n>
