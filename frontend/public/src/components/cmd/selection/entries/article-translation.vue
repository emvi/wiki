<template>
    <emvi-cmd-selection-result :index="index" icon="globe" v-on:click="open" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <div class="item">
                <p>{{entity.name}}</p>
            </div>
            <emvi-cmd-shortcut shortcut="Enter" v-show="active && !isActive">
                {{$t("shortcut_view")}}
            </emvi-cmd-shortcut>
            <span class="label" v-if="isActive">
                {{$t("active")}}
            </span>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {slugWithId} from "../../../../util";
    import {SelectionMixin} from "../mixin";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdEntity from "../../content/entity.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionToggle from "../form/toggle.vue";
    import emviCmdSelectionButton from "../form/button.vue";
    import emviCmdInput from "../../form/input.vue";

    export default {
        mixins: [SelectionMixin],
        components: {
            emviCmdSelectionResult,
            emviCmdEntity,
            emviCmdShortcut,
            emviCmdSelectionToggle,
            emviCmdSelectionButton,
            emviCmdInput
        },
        computed: {
            isActive() {
                return this.entity.id === this.$store.state.page.meta.get("langId");
            }
        },
        watch: {
            enter(enter) {
                if(enter && this.active) {
                    this.open();
                }
            }
        },
        methods: {
            open() {
                let langId = this.$store.state.page.meta.get("langId");

                if(langId === this.entity.id) {
                    return;
                }

                let id = this.$store.state.page.meta.get("id");
                let version = this.$store.state.page.meta.get("version");
                let title = this.$store.state.page.meta.get("title");
                this.$store.dispatch("resetCmd");

                if(version) {
                    this.$router.push(`/read/${slugWithId(title, id)}?version=${version}&lang=${this.entity.id}`);
                }
                else {
                    this.$router.push(`/read/${slugWithId(title, id)}?lang=${this.entity.id}`);
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "active": "Active",
            "shortcut_view": "Open translation"
        },
        "de": {
            "active": "Aktiv",
            "shortcut_view": "Übersetzung öffnen"
        }
    }
</i18n>
