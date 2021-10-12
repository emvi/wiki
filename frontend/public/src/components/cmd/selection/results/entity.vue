<template>
    <emvi-cmd-selection-result cursor-pointer="true" :entity="entity" :index="index" :disabled="disabled" :icon="icon" v-on:click="click" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <emvi-cmd-entity :entity="entity"></emvi-cmd-entity>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-on:click="remove" v-if="!disableRemove" v-show="active">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut :shortcut="$t('key_add')" icon="add" v-on:click="add" v-if="!disableAdd" v-show="active">
                {{$t("shortcut_add")}}
            </emvi-cmd-shortcut>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {PropsMixin} from "../mixin";
    import emviCmdEntity from "../../content/entity.vue";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";

    export default {
        mixins: [PropsMixin],
        components: {
            emviCmdEntity,
            emviCmdSelectionResult,
            emviCmdShortcut
        },
        props: {
            icon: {default: ""},
            disableAdd: {default: false},
            disableRemove: {default: false}
        },
        watch: {
            enter(enter) {
                if(enter && this.active) {
                    this.add();
                }
            },
            del(del) {
                if(del && this.active) {
                    this.remove();
                }
            }
        },
        methods: {
            add() {
                if(!this.disableAdd) {
                    this.$emit("add", this.entity);
                }
            },
            remove() {
                if(!this.disableRemove) {
                    this.$emit("remove", this.entity.id);
                }
            },
            click() {
                if(this.disableAdd) {
                    this.remove();
                }
                else {
                    this.add();
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "key_add": "Enter",
            "key_delete": "Del",
            "shortcut_add": "Add",
            "shortcut_remove": "Remove"
        },
        "de": {
            "key_add": "Enter",
            "key_delete": "Entf",
            "shortcut_add": "Hinzufügen",
            "shortcut_remove": "Entfernen"
        }
    }
</i18n>
