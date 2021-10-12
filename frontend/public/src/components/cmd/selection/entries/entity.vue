<template>
    <emvi-cmd-selection-result :index="index" :disabled="disabled" :icon="icon" :details="details" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <emvi-cmd-entity :entity="entity"></emvi-cmd-entity>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-show="active && !self" v-on:click="toggleDetails(!details)">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
        </template>
        <template slot="details">
            <p>{{removeConfirmation}}</p>
            <emvi-cmd-selection-button icon="back"
                :label="$t('label_no')"
                :index="index"
                :selection="selection"
                :selection-index="0"
                v-on:click="cancel"
                v-on:select="setSelection"></emvi-cmd-selection-button>
            <emvi-cmd-selection-button icon="trash"
                color="red"
                :label="removeLabelYes"
                :index="index"
                :selection="selection"
                :selection-index="1"
                v-on:click="action"
                v-on:select="setSelection"></emvi-cmd-selection-button>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {mapGetters} from "vuex";
    import {SelectionMixin} from "../mixin";
    import emviCmdEntity from "../../content/entity.vue";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionButton from "../form/button.vue";

    export default {
        mixins: [SelectionMixin],
        components: {
            emviCmdEntity,
            emviCmdSelectionResult,
            emviCmdShortcut,
            emviCmdSelectionButton
        },
        props: ["icon", "removeConfirmation", "removeLabelYes"],
        data() {
            return {
                maxSelectionIndex: 1
            };
        },
        computed: {
            ...mapGetters(["user"]),
            self() {
                return this.entity.type === "user" && this.entity.id === this.user.id;
            }
        },
        watch: {
            enter(enter) {
                if(enter && this.active && this.details) {
                    if(this.selection === 0) {
                        this.cancel();
                    } else {
                        this.action();
                    }
                }
            },
            del(del) {
                if(del && this.active && !this.self) {
                    this.toggleDetails(!this.details);
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            }
        },
        methods: {
            action() {
                this.$emit("remove");
            },
            cancel() {
                this.toggleDetails(false);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "key_delete": "Del",
            "shortcut_remove": "Remove",
            "label_no": "No"
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_remove": "Entfernen",
            "label_no": "Nein"
        }
    }
</i18n>
