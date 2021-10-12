<template>
    <div>
        <h5>{{$t("title_profile")}}</h5>
        <emvi-cmd-input :label="$t('label_info')"
            :index="0"
            :error="validation['info']"
            v-model="info"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-input>
            <h5>{{$t("title_notification")}}</h5>
        <emvi-cmd-select :label="$t('label_notifications')"
            :index="1"
            :error="validation['send_notifications_interval']"
            :options="notificationIntervalOptions"
            v-model="sendNotificationsInterval"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-select>
        <emvi-cmd-checkbox :label="$t('label_recommendation_mail')"
            name="recommendation_mail"
            :index="2"
            v-model="recommendationMail"
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
                notificationIntervalOptions: [
                    {value: 0, label: this.$t("select_no_notifications")},
                    {value: 1, label: this.$t("select_notifications_1")},
                    {value: 7, label: this.$t("select_notifications_7")},
                    {value: 30, label: this.$t("select_notifications_30")}
                ],
                info: "",
                sendNotificationsInterval: 0,
                recommendationMail: false
            };
        },
        computed: {
            ...mapGetters(["row", "cmdMeta"])
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
                UserService.getMember()
                    .then(member => {
                        this.info = member.info;
                        this.sendNotificationsInterval = member.send_notifications_interval;
                        this.recommendationMail = member.recommendation_mail;
                    });
            },
            run() {
                this.resetError();
                let data = {
                    info: this.info,
                    send_notifications_interval: parseInt(this.sendNotificationsInterval),
                    recommendation_mail: this.recommendationMail
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
            "title_profile": "Profile",
            "title_notification": "Notification",
            "label_info": "Information",
            "label_notifications": "Unread notifications email",
            "label_recommendation_mail": "Receive article recommendations by email",
            "label_action": "Save",
            "select_no_notifications": "Never",
            "select_notifications_1": "Daily",
            "select_notifications_7": "Once a Week",
            "select_notifications_30": "Once a Month",
            "toast_saved": "Saved."
        },
        "de": {
            "title_profile": "Profil",
            "title_notification": "Benachrichtigung",
            "label_info": "Kurzinfo",
            "label_notifications": "Erhalte ungelesene Benachrichtungen per E-Mail",
            "label_recommendation_mail": "Erhalte Arikelempfehlungen per E-Mail",
            "label_action": "Speichern",
            "select_no_notifications": "Nie",
            "select_notifications_1": "Täglich",
            "select_notifications_7": "Wöchentlich",
            "select_notifications_30": "Monatlich",
            "toast_saved": "Gespeichert."
        }
    }
</i18n>
