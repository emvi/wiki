<template>
    <div>
        <emvi-cmd-submenu icon="user" view="administration-members" :enter="enter" :index="0">
            {{$t("submenu_members")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="mail" view="administration-invitations" :enter="enter" :index="1">
            {{$t("submenu_invitations")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="link" view="administration-invitation-link" :enter="enter" :index="2">
            {{$t("submenu_invitation_link")}}
        </emvi-cmd-submenu>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import emviCmdSubmenu from "../../content/submenu.vue";

    export default {
        components: {emviCmdSubmenu},
        props: ["enter", "tab", "esc", "up", "down"],
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
            },
            tab(tab) {
                if(tab) {
                    if(!tab.shiftKey) {
                        this.$store.dispatch("selectNextRow");
                    }
                    else {
                        this.$store.dispatch("selectPreviousRow");
                    }
                }
            },
            esc(esc) {
                if(esc) {
                    this.$store.dispatch("popColumn");
                }
            },
            up(up) {
                if(up) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down) {
                    this.$store.dispatch("selectNextRow");
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "submenu_members": "Manage Members",
            "submenu_invitations": "Open Invitations",
            "submenu_invitation_link": "Invitation Link"
        },
        "de": {
            "submenu_members": "Mitglieder verwalten",
            "submenu_invitations": "Offene Einladungen",
            "submenu_invitation_link": "Einladungslink"
        }
    }
</i18n>
