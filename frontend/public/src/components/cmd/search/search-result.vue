<template>
    <div v-on:click.stop.prevent="$emit('click')" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false" ref="result">
        <div :class="resultClass">
            <emvi-cmd-avatar :entity="entity" :icon="icon"></emvi-cmd-avatar>
            <slot></slot>
            <emvi-cmd-shortcut shortcut="Tab" icon="chevron" :rotate="showPreview" v-on:click="$emit('preview', index)" v-show="isActive">
                {{$t("preview")}}
            </emvi-cmd-shortcut>
        </div>
        <div v-if="showPreview">
            <emvi-cmd-preview :entity="entity"
                              :up="up"
                              :down="down"
                              v-on:loaded="scroll"></emvi-cmd-preview>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {scrollIntoViewArea, scrollToTopArea} from "../../../util";
    import emviCmdPreview from "./preview.vue";
    import emviCmdAvatar from "../content/avatar.vue";
    import emviCmdShortcut from "../content/shortcut.vue";

    export default {
        components: {
            emviCmdPreview,
            emviCmdAvatar,
            emviCmdShortcut
        },
        props: {
            entity: {},
            icon: "",
            index: {default: 0},
            preview: {default: false},
            tab: {default: false},
            up: {default: false},
            down: {default: false}
        },
        data() {
            return {
                hover: false
            };
        },
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.index === this.row;
            },
            isActive() {
                return this.active || this.hover || this.isTouch;
            },
            showPreview() {
                return this.active && this.preview;
            },
            resultClass() {
                return {
                    "entry cursor-pointer": true,
                    active: this.isActive && !this.isTouch,
                    disabled: this.entity.archived,
                    double: this.isMobile
                };
            }
        },
        watch: {
            row() {
                if(this.active) {
                    scrollIntoViewArea(this.$refs.result, document.getElementById("cmd-results"));
                }
            }
        },
        methods: {
            scroll() {
                this.$nextTick(() => {
                    scrollToTopArea(this.$refs.result, document.getElementById("cmd-results"));
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "preview": "Preview"
        },
        "de": {
            "preview": "Vorschau"
        }
    }
</i18n>
