<template>
    <emvi-cmd-selection-result :index="index" :details="details" :entity="entity" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <emvi-cmd-entity :entity="entity"></emvi-cmd-entity>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-show="active && !self" v-on:click="showDelete">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut shortcut="Tab" icon="key" v-show="active && !self" v-on:click="showDetails">
                {{$t("shortcut_permissions")}}
            </emvi-cmd-shortcut>
        </template>
        <template slot="details">
            <div v-show="permissionsActive">
                <emvi-cmd-selection-toggle name="mod"
                    :label="$t('label_mod')"
                    v-model="mod"
                    :index="index"
                    :selection="selection"
                    :selection-index="0"
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
                    v-on:click="cancel"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
                <emvi-cmd-selection-button icon="trash"
                    color="red"
                    :label="$t('label_yes')"
                    :index="index"
                    :selection="selection"
                    :selection-index="1"
                    v-on:click="remove"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
            </div>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {mapGetters} from "vuex";
    import {SelectionMixin} from "../mixin";
    import {UsergroupService} from "../../../../service";
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
                permissionsActive: false,
                removeActive: false,
                maxSelectionIndex: 1,
                mod: this.entity.is_moderator
            };
        },
        computed: {
            ...mapGetters(["user"]),
            self() {
                return this.entity.id === this.user.id;
            },
            picture() {
                console.log(this.entity);
                return this.entity.picture;
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
                if(del && this.active && !this.self && !this.permissionsActive) {
                    this.showDelete();
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            },
            details(details) {
                if(!details) {
                    this.permissionsActive = false;
                    this.removeActive = false;
                }
            },
            mod() {
                let groupId = this.$store.state.page.meta.get("id");

                UsergroupService.toggleModerator(groupId, this.entity.member_id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            }
        },
        methods: {
            showDetails() {
                this.maxSelectionIndex = 0;
                this.permissionsActive = !this.permissionsActive;
                this.removeActive = false;
                this.toggleDetails(this.permissionsActive);
            },
            showDelete() {
                this.maxSelectionIndex = 1;
                this.permissionsActive = false;
                this.removeActive = !this.removeActive;
                this.toggleDetails(this.removeActive);
            },
            remove() {
                this.resetError();
                let groupId = this.$store.state.page.meta.get("id");

                UsergroupService.removeMember(groupId, this.entity.member_id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_removed"));
                        this.$emit("remove", this.entity.id);
                        this.cancel();
                        this.$store.dispatch("setMeta", {key: "updateList", value: true});
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.permissionsActive = false;
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
            "shortcut_permissions": "Permissions",
            "shortcut_remove": "Remove",
            "confirmation": "Are you sure you would like to remove the member from this group?",
            "label_mod": "Moderator",
            "label_no": "No",
            "label_yes": "Yes, remove member",
            "toast_saved": "Saved.",
            "toast_removed": "The member has been removed from the group."
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_permissions": "Berechtigungen",
            "shortcut_remove": "Entfernen",
            "confirmation": "Möchtest du dieses Mitglied wirklich aus der Gruppe entfernen?",
            "label_mod": "Moderator",
            "label_no": "Nein",
            "label_yes": "Ja, Mitglied entfernen",
            "toast_saved": "Gespeichert.",
            "toast_removed": "Das Mitglied wurde aus der Gruppe entfernt."
        }
    }
</i18n>
