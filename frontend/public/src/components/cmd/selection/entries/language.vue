<template>
    <emvi-cmd-selection-result :index="index" icon="globe" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <div class="item">
                <p>{{entity.name}} ({{entity.code}})</p>
                <small v-if="entity.default">{{$t("label_default")}}</small>
            </div>
            <emvi-cmd-shortcut shortcut="Tab" icon="check" v-show="active" v-on:click="setDefault">
                {{$t("shortcut_default")}}
            </emvi-cmd-shortcut>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {SelectionMixin} from "../mixin";
    import {LangService} from "../../../../service";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";

    export default {
        mixins: [SelectionMixin],
        components: {emviCmdSelectionResult, emviCmdShortcut},
        watch: {
            tab(tab) {
                if(tab && this.active) {
                    this.setDefault();
                }
            }
        },
        methods: {
            setDefault() {
                this.resetError();
                LangService.switchDefault(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                        this.$emit("update");
                        this.focusCmdInput();
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label_default": "Default",
            "shortcut_default": "Set as default",
            "toast_saved": "Saved."
        },
        "de": {
            "label_default": "Standard",
            "shortcut_default": "Als Standard setzen",
            "toast_saved": "Gespeichert."
        }
    }
</i18n>
