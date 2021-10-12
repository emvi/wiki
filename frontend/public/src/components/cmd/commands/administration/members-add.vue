<template>
    <div>
        <emvi-cmd-input :label="$t('label_email')"
            :index="0"
            :hint="$t('hint_email')"
            :error="err"
            v-model="email"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="add"
            v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-selection-email v-for="(email, index) in emails"
            :key="email"
            :entity="email"
            :index="index+1"
            :tab="tab"
            :del="del"
            :up="up"
            :down="down"
            v-on:remove="remove"></emvi-cmd-selection-email>
        <emvi-cmd-checkbox :label="$t('label_readonly')"
            :index="emails.length+1"
            :hint="$t('hint_readonly')"
            v-model="readOnly"
            name="readonly"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="invite"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button :icon="icon"
            :label="$t('label_invite')"
            :index="emails.length+2"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="invite"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import ISO6391 from "iso-639-1";
    import {updateSelectedRow} from "../../util";
    import {MemberService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";
    import emviCmdSelectionEmail from "../../selection/results/email.vue";

    export default {
        components: {emviCmdInput, emviCmdCheckbox, emviCmdButton, emviCmdSelectionEmail},
        props: ["up", "down", "enter", "tab", "del", "esc"],
        data() {
            return {
                icon: "send",
                err: "",
                email: "",
                readOnly: false,
                emails: []
            };
        },
        computed: {
            ...mapGetters(["row"]),
            languageOptions() {
                let langs = ISO6391.getLanguages(ISO6391.getAllCodes());
                let options = [];

                for(let i = 0; i < langs.length; i++) {
                    options.push({value: langs[i].code, label: langs[i].nativeName});
                }

                return options;
            }
        },
        watch: {
            email() {
                if(this.err) {
                    this.err = "";
                }
            },
            row(row) {
                updateSelectedRow(row, this.emails.length+3, this.$store);
            },
            enter(enter) {
                if(enter && this.row > 0) {
                    this.invite();
                }
            },
            tab(tab) {
                if(tab && this.isEmailRow()) {
                    if(tab.shiftKey) {
                        this.$store.dispatch("selectPreviousRow");
                    }
                    else {
                        this.$store.dispatch("selectNextRow");
                    }
                }
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            },
            up(up) {
                if(up && this.isEmailRow()) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down && this.isEmailRow()) {
                    this.$store.dispatch("selectNextRow");
                }
            }
        },
        methods: {
            add() {
                this.email = this.email.trim();
                let match = /^\S+@\S+$/;

                if(this.email.length && match.test(this.email)) {
                    this.emails.push(this.email);
                    this.email = "";
                }
                else {
                    this.err = this.$t("error_email");
                }
            },
            remove(email) {
                let index = this.emails.indexOf(email);

                if(index > -1) {
                    this.emails.splice(index, 1);
                    this.$store.dispatch("selectRow", this.row-1);
                }
            },
            invite() {
                this.resetError();
                this.icon = "sync";

                MemberService.inviteMember(this.emails, this.readOnly)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_invited"));
                        this.$store.dispatch("popColumn");
                    })
                    .catch(e => {
                        this.icon = "send";
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
            },
            isEmailRow() {
                return this.row > 0 && this.row < this.emails.length+2;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label_email": "Email",
            "label_readonly": "Read only",
            "label_invite": "Send invitations",
            "hint_email": "Press ENTER to add to invitation list.",
            "hint_readonly": "If set, the invited Members will have read access only. This can be changed later.",
            "toast_invited": "Invitations have been sent.",
            "error_email": "The email address is invalid."
        },
        "de": {
            "label_email": "E-Mail",
            "label_readonly": "Nur Lesezugriff",
            "label_invite": "Einladungen senden",
            "hint_email": "ENTER um zur Einladungsliste hinzuzufügen.",
            "hint_readonly": "Wenn gesetzt, erhalten die eingeladenen Mitglieder ausschließlich Lesezugriff. Das kann später geändert werden.",
            "toast_invited": "Die Einladungen wurden verschickt.",
            "error_email": "Die E-Mail-Adresse ist ungültig."
        }
    }
</i18n>
