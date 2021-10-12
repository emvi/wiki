<template>
    <div>
        <emvi-cmd-radio :label="$t('label_groups')"
                        name="groups"
                        :index="0"
                        :disabled="!isExpert"
                        v-model="createGroup"
                        :options="createGroupOptions"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-radio>
        <emvi-cmd-button icon="save"
                         color="green"
                         :label="isExpert ? $t('label_save') : $t('label_save')+' '+$t('expert')"
                         :index="1"
                         :disabled="!isExpert"
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
    import emviCmdRadio from "../../form/radio.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdRadio, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                createGroupOptions: [
                    {label: this.$t("label_permissions_groups_admins"), value: 0},
                    {label: this.$t("label_permissions_groups_mods"), value: 1},
                    {label: this.$t("label_permissions_groups_everyone"), value: 2}
                ],
                createGroup: 2
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
                    this.$store.dispatch("popColumn");
                }
            }
        },
        beforeMount() {
            if(this.organization.create_group_mod) {
                this.createGroup = 1;
            }
            else if(this.organization.create_group_admin) {
                this.createGroup = 0;
            }
        },
        methods: {
            save() {
                if(!this.isExpert) {
                    return;
                }

                this.resetError();
                this.createGroup = parseInt(this.createGroup);
                let mod = this.createGroup === 1;
                let admin = this.createGroup === 0 || mod;

                OrganizationService.updateOrganizationPermissions(admin, mod)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                        this.$store.dispatch("loadOrganization");
                        this.$store.dispatch("popColumn");
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
            "label_groups": "Groups",
            "label_save": "Save",
            "toast_saved": "Saved.",
            "label_permissions_groups_admins": "Only Administrators can create and manage groups",
            "label_permissions_groups_mods": "Administrators and moderators can create and manage groups",
            "label_permissions_groups_everyone": "Every user can create and manage groups",
            "expert": "(requires an Expert organization)"
        },
        "de": {
            "label_groups": "Gruppen",
            "label_save": "Speichern",
            "toast_saved": "Gespeichert.",
            "label_permissions_groups_admins": "Nur Administratoren können Gruppen anlegen und verwalten",
            "label_permissions_groups_mods": "Administratoren und Moderatoren können Gruppen anlegen und verwalten",
            "label_permissions_groups_everyone": "Jeder Nutzer kann Gruppen anlegen und verwalten",
            "expert": "(benötigt eine Expert Organisation)"
        }
    }
</i18n>
