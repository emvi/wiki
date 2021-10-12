<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-input :label="$t('label_name')"
                        :index="0"
                        :hint="organization.name"
                        :error="validation['name']"
                        v-model="name"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-button icon="trash"
                         color="red"
                         :label="$t('label_delete')"
                         :index="1"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="save"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {OrganizationService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdInput, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                name: ""
            };
        },
        computed: {
            ...mapGetters(["row", "organization"])
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
        methods: {
            save() {
                this.resetError();

                OrganizationService.deleteOrganization(this.name)
                    .then(() => {
                        this.$store.dispatch("redirectOrganizations");
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
            "text": "Are you sure you want to delete this organization? This step cannot be undone!",
            "label_name": "Name of this Organization",
            "label_delete": "Delete Organization"
        },
        "de": {
            "text": "Soll diese Organisation wirklich gelöscht werden? Die Löschung kann nicht rückgängig gemacht werden!",
            "label_name": "Name der Organisation",
            "label_delete": "Organisation löschen"
        }
    }
</i18n>
