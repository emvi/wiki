<template>
    <div>
        <emvi-cmd-submenu icon="user" view="settings-account" :enter="enter" :index="0">
            {{$t("submenu_account")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="infobox" view="settings-ui" :enter="enter" :index="1">
            {{$t("submenu_ui")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="help" view="introduction" :enter="enter" :index="2">
            {{$t("submenu_introduction")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="back" color="red" view="settings-leave" :enter="enter" :index="3">
            {{$t("submenu_leave")}}
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
                updateSelectedRow(row, 4, this.$store);
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
            "submenu_account": "Profile and Notification Settings",
            "submenu_ui": "User Interface",
            "submenu_introduction": "Repeat Introduction",
            "submenu_leave": "Leave Organization"
        },
        "de": {
            "submenu_account": "Profil- und Benachrichtungseinstellungen",
            "submenu_ui": "Benutzeroberfläche",
            "submenu_introduction": "Einführung wiederholen",
            "submenu_leave": "Organisation verlassen"
        }
    }
</i18n>
