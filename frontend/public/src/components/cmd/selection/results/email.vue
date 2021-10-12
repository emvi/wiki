<template>
    <emvi-cmd-selection-result :index="index" icon="mail" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <div class="item">
                <p>{{entity}}</p>
            </div>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-on:click="showDelete" v-show="active">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {SelectionMixin} from "../mixin";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";

    export default {
        mixins: [SelectionMixin],
        components: {emviCmdSelectionResult, emviCmdShortcut},
        watch: {
            active(active) {
                if(active) {
                    this.focusCmdInput();
                }
            },
            del(e) {
                if(e && this.active) {
                    e.preventDefault();
                    e.stopPropagation();
                    this.showDelete();
                }
            }
        },
        methods: {
            showDelete() {
                this.$emit("remove", this.entity);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "key_delete": "Del",
            "shortcut_remove": "Remove"
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_remove": "Entfernen"
        }
    }
</i18n>
