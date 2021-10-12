<template>
    <div>
        <h5>{{$t("title")}}</h5>
        <emvi-cmd-checkbox :label="$t('label_show_create_button')"
            name="show_create_button"
            :index="0"
            v-model="showCreateButton"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_show_navigation')"
            name="show_navigation"
            :index="1"
            v-model="showNavigation"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_show_action_buttons')"
            name="show_action_buttons"
            :index="2"
            v-model="showActionButtons"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="save"
            color="green"
            :label="$t('label_action')"
            :index="3"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {UserService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdSelect from "../../form/select.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";
    import {isEmptyObject} from "../../../../util";

    export default {
        components: {emviCmdInput, emviCmdSelect, emviCmdCheckbox, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                showCreateButton: false,
                showNavigation: false,
                showActionButtons: false
            };
        },
        computed: {
            ...mapGetters(["row", "cmdMeta", "member"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 4, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        mounted() {
            this.loadSettings();
        },
        methods: {
            loadSettings() {
                this.showCreateButton = this.member.show_create_button;
                this.showNavigation = this.member.show_navigation;
                this.showActionButtons = this.member.show_action_buttons;
            },
            run() {
                this.resetError();
                let data = {
                    show_create_button: this.showCreateButton,
                    show_navigation: this.showNavigation,
                    show_action_buttons: this.showActionButtons
                };

                UserService.saveMember(data)
                    .then(() => {
                        this.$store.dispatch("reload");
                        this.$store.dispatch("success", this.$t("toast_saved"));
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
            "title": "User Interface",
            "label_show_create_button": "Show create button in menu",
            "label_show_navigation": "Show navigation",
            "label_show_action_buttons": "Show buttons for actions on page",
            "label_action": "Save",
            "toast_saved": "Saved."
        },
        "de": {
            "title": "Benutzeroberfläche",
            "label_show_create_button": "Zeige Erstellen-Button im Menü",
            "label_show_navigation": "Zeige Navigation",
            "label_show_action_buttons": "Zeige Buttons für Aktionen auf der Seite",
            "label_action": "Speichern",
            "toast_saved": "Gespeichert."
        }
    }
</i18n>
