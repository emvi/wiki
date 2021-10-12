<template>
    <div>
        <emvi-cmd-input :label="$t('label_name')"
                        :index="0"
                        :error="validation['name']"
                        v-model="name"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-input :label="$t('label_domain')"
                        :index="1"
                        :error="validation['domain']"
                        :hint="$t('note_update_orga')"
                        v-model="domain"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-button icon="save"
                         color="green"
                         :label="$t('label_save')"
                         :index="2"
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
                name: "",
                domain: ""
            };
        },
        computed: {
            ...mapGetters(["row", "organization"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.$store.dispatch("popColumn");
                }
            }
        },
        mounted() {
            this.name = this.organization.name;
            this.domain = this.organization.name_normalized;
        },
        methods: {
            save() {
                this.resetError();

                if(this.name !== this.organization.name || this.domain !== this.organization.name_normalized) {
                    OrganizationService.updateOrganization(this.name, this.domain)
                        .then(() => {
                            if(this.domain === this.organization.name_normalized) {
                                this.$store.dispatch("success", this.$t("toast_saved"));
                                this.$store.dispatch("loadOrganization");
                                this.$store.dispatch("popColumn");
                            }
                            else {
                                this.$store.dispatch("redirectOrganizations");
                            }
                        })
                        .catch(e => {
                            this.setError(e);
                        });
                }
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
            "label_name": "Organization Name",
            "label_domain": "Subdomain",
            "label_save": "Save",
            "toast_saved": "Saved.",
            "note_update_orga": "The subdomain is the first part of the domain (my-orga.emvi.com for example). If you change the URL, links containing your original URL will NOT be forwarded! You will be redirected to the organizations page, when you change the name or URL."
        },
        "de": {
            "label_name": "Name der Organisation",
            "label_domain": "Subdomain",
            "label_save": "Speichern",
            "toast_saved": "Gespeichert.",
            "note_update_orga": "Die Subdomain ist der erste Teil der Domain (z.B. meine-orga.emvi.com). Wenn du die URL änderst, werden Links mit der ursprünglichen URL nicht weitergeleitet. Du wirst zur Organisationsübersicht umgeleitet, nachdem du Name oder URL geändert hast."
        }
    }
</i18n>
