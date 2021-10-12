<template>
    <div>
        <emvi-cmd-radio :label="$t('label_access')"
            :hint="$t('hint_access')"
            :options="accessModeOptions"
            :index="0"
            v-model="accessMode"
            name="accessMode"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:esc="cancel"></emvi-cmd-radio>
        <emvi-cmd-checkbox :label="isExpert ? $t('label_external') : $t('label_external')+' '+$t('expert')"
            :hint="$t('hint_external')"
            :index="1"
            :disabled="!isExpert"
            v-model="clientAccess"
            name="clientAccess"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-submenu v-show="showMemberSubmenu"
            icon="user"
            view="article-permissions-member"
            :index="2"
            :enter="enter"
            :tab="tab"
            v-on:next="nextRow"
            v-on:previous="previousRow">
            {{$t("submenu_member")}}
        </emvi-cmd-submenu>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import emviCmdRadio from "../../form/radio.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdSubmenu from "../../content/submenu.vue";

    export default {
        components: {emviCmdRadio, emviCmdCheckbox, emviCmdSubmenu},
        props: ["enter", "tab", "esc", "up", "down"],
        data() {
            return {
                accessModeOptions: [
                    {label: this.$t("accessmode_public_write"), value: "0"},
                    {label: this.$t("accessmode_public_read"), value: "1"},
                    {label: this.$t("accessmode_limited"), value: "2"},
                    {label: this.$t("accessmode_private"), value: "3"}
                ],
                accessMode: String(this.$store.state.page.meta.get("accessMode")),
                clientAccess: this.$store.state.page.meta.get("clientAccess")
            };
        },
        computed: {
            ...mapGetters(["row", "metaUpdate"]),
            showMemberSubmenu() {
                let mode = parseInt(this.accessMode);
                return mode === 1 || mode === 2;
            }
        },
        watch: {
            row(row) {
                updateSelectedRow(row, this.showMemberSubmenu ? 3 : 2, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            },
            up(up) {
                if(up && this.row === 2) {
                    this.previousRow();
                }
            },
            down(down) {
                if(down && this.row === 2) {
                    this.nextRow();
                }
            },
            metaUpdate() {
                let accessMode = this.$store.state.page.meta.get("accessMode");
                let clientAccess = this.$store.state.page.meta.get("clientAccess");

                if(accessMode !== undefined && accessMode !== parseInt(this.accessMode)) {
                    this.accessMode = accessMode;
                }

                if(clientAccess !== undefined && clientAccess !== this.clientAccess) {
                    this.clientAccess = clientAccess;
                }
            },
            accessMode() {
                if(this.clientAccess && this.accessMode > 1) {
                    this.clientAccess = false;
                }

                this.$store.dispatch("setMetaVars", [
                    {key: "accessMode", value: parseInt(this.accessMode)},
                    {key: "clientAccess", value: this.clientAccess}
                ]);
            },
            clientAccess() {
                if(this.clientAccess && parseInt(this.accessMode) > 1) {
                    this.accessMode = "1";
                }

                this.$store.dispatch("setMetaVars", [
                    {key: "accessMode", value: parseInt(this.accessMode)},
                    {key: "clientAccess", value: this.clientAccess}
                ]);
            }
        },
        methods: {
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
            "label_access": "Visibility",
            "label_external": "External Access",
            "hint_access": "Manage who inside this organization can read and edit this article.",
            "hint_external": "Allow this article to be read externally through the Client API.",
            "accessmode_public_write": "Full access for everyone",
            "accessmode_public_read": "Read access for everyone, specify who has write access",
            "accessmode_limited": "Set read and write access for individual members",
            "accessmode_private": "Just you have access",
            "submenu_member": "Manage Member Access",
            "expert": "(requires an Expert Organization)"
        },
        "de": {
            "label_access": "Sichtbarkeit",
            "label_external": "Externer Zugriff",
            "hint_access": "Legt fest wer den Artikel innerhalb der Organisation sehen und bearbeiten kann.",
            "hint_external": "Erlaubt den externen Zugriff über die Client API.",
            "accessmode_public_write": "Vollzugriff für jeden",
            "accessmode_public_read": "Lesezugriff für jeden, definiere wer Schreibzugriff hat",
            "accessmode_limited": "Setze Lese- und Schreibzugriff für einzelne Mitglieder",
            "accessmode_private": "Nur du hast Zugriff",
            "submenu_member": "Mitgliederzugriff verwalten",
            "expert": "(benötigt eine Expert Organisation)"
        }
    }
</i18n>
