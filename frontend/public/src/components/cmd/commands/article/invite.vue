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
        <emvi-cmd-button icon="add-user"
            :label="$t('label_send')"
            :index="2"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="send"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {ArticleService} from "../../../../service";
    import emviCmdSubmenu from "../../content/submenu.vue";
    import emviCmdTextarea from "../../form/textarea.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdSubmenu, emviCmdTextarea, emviCmdButton},
        props: ["enter", "tab", "esc", "up", "down"],
        data() {
            return {
                members: [],
                message: ""
            };
        },
        computed: {
            ...mapGetters(["row", "cmdMeta"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
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
                }

                if(this.cmdMeta.get("members")) {
                    this.members = this.cmdMeta.get("members");
                }
            },
            send() {
                this.resetError();
                let id = this.$store.state.page.meta.get("id");
                let langId = this.$store.state.page.meta.get("langId");
                let room = this.$store.state.page.meta.get("room");
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

                ArticleService.inviteArticle(id, langId, room, user, groups, this.message)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_invited"));
                        this.$store.dispatch("resetCmd");
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            selectMembers() {
                this.$store.dispatch("setCmdMeta", {key: "message", value: this.message});
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
            "label_send": "Invite",
            "toast_invited": "Invitations were send."
        },
        "de": {
            "selected": "Mitglieder ausgewählt",
            "label_message": "Nachricht",
            "label_send": "Einladen",
            "toast_invited": "Einladungen wurden versendet."
        }
    }
</i18n>
