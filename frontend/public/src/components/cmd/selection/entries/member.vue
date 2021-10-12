<template>
    <emvi-cmd-selection-result :index="index" :details="details" :entity="entity" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <emvi-cmd-entity :entity="entity"></emvi-cmd-entity>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-on:click="showDelete" v-show="active && !self">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut shortcut="Tab" icon="key" v-on:click="showDetails" v-show="active && !self">
                {{$t("shortcut_roles")}}
            </emvi-cmd-shortcut>
        </template>
        <template slot="details">
            <div v-show="rolesActive">
                <p v-if="!isExpert">{{$t("expert")}}</p>
                <emvi-cmd-selection-toggle name="admin"
                    :label="$t('label_admin')"
                    v-model="admin"
                    :index="index"
                    :disabled="!isExpert"
                    :selection="selection"
                    :selection-index="0"
                    :enter="enter"
                    v-on:select="setSelection"></emvi-cmd-selection-toggle>
                <emvi-cmd-selection-toggle name="mod"
                    :label="$t('label_mod')"
                    v-model="mod"
                    :index="index"
                    :disabled="!isExpert"
                    :selection="selection"
                    :selection-index="1"
                    :enter="enter"
                    v-on:select="setSelection"></emvi-cmd-selection-toggle>
                <emvi-cmd-selection-toggle name="ro"
                    :label="$t('label_ro')"
                    v-model="ro"
                    :index="index"

                    :selection="selection"
                    :selection-index="2"
                    :enter="enter"
                    v-on:select="setSelection"></emvi-cmd-selection-toggle>
            </div>
            <div v-show="removeActive">
                <p>{{$t("confirmation")}}</p>
                <emvi-cmd-selection-button icon="back"
                    :label="$t('label_no')"
                    :index="index"
                    :selection="selection"
                    :selection-index="0"
                    v-on:select="setSelection"
                    v-on:click="cancel"></emvi-cmd-selection-button>
                <emvi-cmd-selection-button icon="trash"
                    color="red"
                    :label="$t('label_yes')"
                    :index="index"
                    :selection="selection"
                    :selection-index="1"
                    v-on:select="setSelection"
                    v-on:click="remove"></emvi-cmd-selection-button>
            </div>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {mapGetters} from "vuex";
    import {SelectionMixin} from "../mixin";
    import {MemberService} from "../../../../service";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdEntity from "../../content/entity.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionToggle from "../form/toggle.vue";
    import emviCmdSelectionButton from "../form/button.vue";

    export default {
        mixins: [SelectionMixin],
        components: {
            emviCmdSelectionResult,
            emviCmdEntity,
            emviCmdShortcut,
            emviCmdSelectionToggle,
            emviCmdSelectionButton
        },
        data() {
            return {
                rolesActive: false,
                removeActive: false,
                maxSelectionIndex: 1,
                admin: this.entity.organization_member.is_admin,
                mod: this.entity.organization_member.is_moderator,
                ro: this.entity.organization_member.read_only
            };
        },
        computed: {
            ...mapGetters(["user"]),
            self() {
                return this.entity.id === this.user.id;
            }
        },
        watch: {
            enter(enter) {
                if(enter && this.active && this.details && this.removeActive) {
                    if(this.selection === 0) {
                        this.cancel();
                    } else {
                        this.remove();
                    }
                }
            },
            tab(tab) {
                if(tab && this.active && !this.self && !this.removeActive) {
                    this.showDetails();
                }
            },
            del(del) {
                if(del && this.active && !this.self && !this.rolesActive) {
                    this.showDelete();
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            },
            admin(admin) {
                if(admin) {
                    this.mod = true;
                    this.ro = false;
                }

                MemberService.toggleAdmin(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            mod(mod) {
                if(mod) {
                    this.ro = false;
                }
                else {
                    this.admin = false;
                }

                MemberService.toggleModerator(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            ro(ro) {
                if(ro) {
                    this.admin = false;
                    this.mod = false;
                }

                MemberService.toggleReadOnly(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            details(details) {
                if(!details) {
                    this.rolesActive = false;
                    this.removeActive = false;
                }
            }
        },
        methods: {
            showDetails() {
                this.maxSelectionIndex = 2;
                this.rolesActive = !this.rolesActive;
                this.removeActive = false;
                this.toggleDetails(this.rolesActive);
            },
            showDelete() {
                this.maxSelectionIndex = 1;
                this.rolesActive = false;
                this.removeActive = !this.removeActive;
                this.toggleDetails(this.removeActive);
            },
            remove() {
                this.resetError();
                MemberService.removeMember(this.entity.id, false)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_removed"));
                        this.$emit("remove", this.entity.id);
                        this.cancel();
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.rolesActive = false;
                this.removeActive = false;
                this.toggleDetails(false);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "key_delete": "Del",
            "shortcut_roles": "Roles",
            "shortcut_remove": "Remove",
            "confirmation": "Are you sure you would like to remove this member?",
            "label_admin": "Administrator",
            "label_mod": "Moderator",
            "label_ro": "Read only",
            "label_no": "No",
            "label_yes": "Yes, remove member",
            "toast_saved": "Saved.",
            "toast_removed": "The member has been removed from the organization.",
            "expert": "To change the role of a user you need an Expert organization."
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_roles": "Rollen",
            "shortcut_remove": "Entfernen",
            "confirmation": "Möchtest du dieses Mitglied wirklich entfernen?",
            "label_admin": "Administrator",
            "label_mod": "Moderator",
            "label_ro": "nur Lesezugriff",
            "label_no": "Nein",
            "label_yes": "Ja, Mitglied entfernen",
            "toast_saved": "Gespeichert.",
            "toast_removed": "Das Mitglied wurde aus der Organisation entfernt.",
            "expert": "Um die Rolle eines Mitglieds zu ändern benötigst du eine Expert Organisation."
        }
    }
</i18n>
