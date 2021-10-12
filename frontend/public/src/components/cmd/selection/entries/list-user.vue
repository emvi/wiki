<template>
    <emvi-cmd-selection-result :index="index" :active="active" :details="details" :entity="entity" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <emvi-cmd-entity :entity="entity"></emvi-cmd-entity>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-on:click="showDelete" v-show="!self">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut shortcut="Tab" icon="key" v-on:click="showDetails" v-show="active && !self">
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
    import {ArticlelistService} from "../../../../service";
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
                let listId = this.$store.state.page.meta.get("id");

                ArticlelistService.toggleModerator(listId, this.entity.member_id)
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
                this.$emit("remove", this.entity.id);
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
            "confirmation": "Are you sure you would like to remove the member from this list?",
            "label_mod": "Moderator",
            "label_no": "No",
            "label_yes": "Yes, remove member",
            "toast_saved": "Saved."
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_permissions": "Berechtigungen",
            "shortcut_remove": "Entfernen",
            "confirmation": "Möchtest du dieses Mitglied wirklich aus der Liste entfernen?",
            "label_mod": "Moderator",
            "label_no": "Nein",
            "label_yes": "Ja, Mitglied entfernen",
            "toast_saved": "Gespeichert."
        }
    }
</i18n>
