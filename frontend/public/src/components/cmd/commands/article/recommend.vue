<template>
    <div>
        <emvi-cmd-submenu icon="group"
                          show-shortcut="true"
                          :index="0"
                          :enter="enter"
                          v-on:enter="selectMembers">
            {{members.length}} {{$t("selected")}}
        </emvi-cmd-submenu>
        <emvi-cmd-textarea :label="$t('label_message')"
                           :index="1"
                           optional="true"
                           v-model="message"
                           v-on:next="nextRow"
                           v-on:previous="previousRow"
                           v-on:enter="send"
                           v-on:esc="cancel"></emvi-cmd-textarea>
        <emvi-cmd-checkbox :label="$t('label_receive_confirmation')"
                           :index="2"
                           name="receiveConfirmation"
                           v-model="receiveConfirmation"
                           v-on:next="nextRow"
                           v-on:previous="previousRow"
                           v-on:enter="send"
                           v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="publish"
                         :label="$t('label_send')"
                         :index="3"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="send"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import emviCmdSubmenu from "../../content/submenu.vue";
    import emviCmdTextarea from "../../form/textarea.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";
    import {ArticleService} from "../../../../service";

    export default {
        components: {emviCmdSubmenu, emviCmdTextarea, emviCmdCheckbox, emviCmdButton},
        props: ["enter", "tab", "esc", "up", "down"],
        data() {
            return {
                members: [],
                message: "",
                receiveConfirmation: false
            };
        },
        computed: {
            ...mapGetters(["row", "cmdMeta"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 4, this.$store);
            },
            tab(tab) {
                if(tab && this.row === 0) {
                    if(!tab.shiftKey) {
                        this.$store.dispatch("selectNextRow");
                    }
                    else {
                        this.$store.dispatch("selectPreviousRow");
                    }
                }
            },
            esc(esc) {
                if(esc && this.row === 0) {
                    this.cancel();
                }
            },
            up(up) {
                if(up && this.row === 0) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down && this.row === 0) {
                    this.$store.dispatch("selectNextRow");
                }
            }
        },
        beforeMount() {
            this.init();
        },
        methods: {
            init() {
                if(this.cmdMeta.get("message")) {
                    this.message = this.cmdMeta.get("message");
                    this.receiveConfirmation = this.cmdMeta.get("receiveConfirmation");
                }

                if(this.cmdMeta.get("members")) {
                    this.members = this.cmdMeta.get("members");
                }
            },
            send() {
                this.resetError();
                let id = this.$store.state.page.meta.get("id");
                let user = [];
                let groups = [];

                for(let i = 0; i < this.members.length; i++) {
                    if(this.members[i].type === "user") {
                        user.push(this.members[i].id);
                    }
                    else {
                        groups.push(this.members[i].id);
                    }
                }

                ArticleService.recommendArticle(id, user, groups, this.message, this.receiveConfirmation)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_recommended"));
                        this.$store.dispatch("resetCmd");
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            selectMembers() {
                this.$store.dispatch("setCmdMetaVars", [
                    {key: "message", value: this.message},
                    {key: "receiveConfirmation", value: this.receiveConfirmation}
                ]);
                this.$store.dispatch("pushColumn", "select-user-group");
                this.$store.dispatch("setCmdMeta", {key: "members", value: this.members});
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
            "selected": "members selected",
            "label_message": "Message",
            "label_receive_confirmation": "Request a read confirmation",
            "label_send": "Send",
            "toast_recommended": "Recommendations were send."
        },
        "de": {
            "selected": "Mitglieder ausgewählt",
            "label_message": "Nachricht",
            "label_receive_confirmation": "Lesebestätigung anfordern",
            "label_send": "Senden",
            "toast_recommended": "Empfehlungen wurden versendet."
        }
    }
</i18n>
