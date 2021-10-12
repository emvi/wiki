<template>
    <div>
        <emvi-cmd-organization-header></emvi-cmd-organization-header>
        <emvi-cmd-submenu icon="img" view="administration-image" :enter="enter" :index="0">
            {{$t("submenu_image")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="settings" view="administration-general" :enter="enter" :index="1">
            {{$t("submenu_general")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="key" view="administration-permission" :enter="enter" :index="2">
            {{$t("submenu_permission")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="billing" :enter="enter" :index="3" v-on:enter="openBillingPage">
            {{$t("submenu_billing")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="user" view="administration-members-invitations" :enter="enter" :index="4">
            {{$t("submenu_members_invitations")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="globe" view="administration-languages" :enter="enter" :index="5">
            {{$t("submenu_languages")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="command" view="administration-clients" :enter="enter" :index="6">
            {{$t("submenu_clients")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="error" color="red" view="administration-delete" :enter="enter" :index="7">
            {{$t("submenu_delete")}}
        </emvi-cmd-submenu>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import emviCmdOrganizationHeader from "../../content/organization-header.vue";
    import emviCmdSubmenu from "../../content/submenu.vue";

    export default {
        components: {emviCmdOrganizationHeader, emviCmdSubmenu},
        props: ["enter", "tab", "esc", "up", "down"],
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 8, this.$store);
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
        },
        methods: {
            openBillingPage() {
                this.$router.push("/billing");
                this.$store.dispatch("resetCmd");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "submenu_image": "Picture",
            "submenu_general": "General settings",
            "submenu_permission": "Permissions",
            "submenu_billing": "Billing",
            "submenu_members_invitations": "Members and invitations",
            "submenu_languages": "Languages",
            "submenu_clients": "Clients",
            "submenu_delete": "Delete Organization"
        },
        "de": {
            "submenu_image": "Bild",
            "submenu_general": "Allgemeine Einstellungen",
            "submenu_permission": "Berechtigungen",
            "submenu_billing": "Abrechnung",
            "submenu_members_invitations": "Mitglieder und Einladungen",
            "submenu_languages": "Sprachen",
            "submenu_clients": "Clients",
            "submenu_delete": "Organisation löschen"
        }
    }
</i18n>
