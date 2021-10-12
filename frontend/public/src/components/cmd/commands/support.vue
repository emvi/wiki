<template>
    <div>
        <emvi-cmd-select :label="$t('label_type')"
            :index="0"
            :options="typeOptions"
            v-model="type"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-select>
        <emvi-cmd-input :label="$t('label_subject')"
            :index="1"
            :error="validation['subject']"
            v-model="subject"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-textarea :label="$t('label_message')"
            :index="2"
            :error="validation['message']"
            v-model="message"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-textarea>
        <emvi-cmd-button icon="send"
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
    import {updateSelectedRow} from "../util";
    import {SupportService} from "../../../service";
    import emviCmdSelect from "../form/select.vue";
    import emviCmdInput from "../form/input.vue";
    import emviCmdTextarea from "../form/textarea.vue";
    import emviCmdButton from "../form/button.vue";

    export default {
        components: {emviCmdSelect, emviCmdInput, emviCmdTextarea, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                typeOptions: [
                    {value: "type_question", label: this.$t("type_question")},
                    {value: "type_feature", label: this.$t("type_feature")},
                    {value: "type_feedback", label: this.$t("type_feedback")},
                    {value: "type_issue", label: this.$t("type_issue")},
                    {value: "type_billing", label: this.$t("type_billing")}
                ],
                type: "type_question",
                subject: "",
                message: ""
            };
        },
        computed: {
            ...mapGetters(["row"])
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
        methods: {
            run() {
                this.resetError();

                SupportService.contact(this.type, this.subject, this.message)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_send"));
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
            "label_type": "What is the topic of your request?",
            "label_subject": "Subject",
            "label_message": "Message",
            "label_action": "Send",
            "type_question": "General Question",
            "type_feature": "Feature Request",
            "type_feedback": "Feedback",
            "type_issue": "Technical Issue",
            "type_billing": "Billing",
            "toast_send": "Your ticket has been created."
        },
        "de": {
            "label_type": "Was ist dein Anliegen?",
            "label_subject": "Betreff",
            "label_message": "Nachricht",
            "label_action": "Absenden",
            "type_question": "Generelle Frage",
            "type_feature": "Funktionalität vorschlagen",
            "type_feedback": "Rückmeldung zum Produkt",
            "type_issue": "Technisches Problem",
            "type_billing": "Abrechnung",
            "toast_send": "Dein Ticket wurde erstellt."
        }
    }
</i18n>
