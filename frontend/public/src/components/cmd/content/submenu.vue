<template>
    <div :class="colorClass" v-on:click.stop.prevent="open">
        <div :class="submenuClass" ref="submenu">
            <emvi-cmd-icon :icon="icon"></emvi-cmd-icon>
            <div class="item">
                <p>
                    <slot></slot>
                </p>
            </div>
            <emvi-cmd-shortcut shortcut="Enter" v-if="showShortcut" v-show="active"></emvi-cmd-shortcut>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {scrollIntoViewArea} from "../../../util";
    import emviCmdIcon from "./icon.vue";
    import emviCmdShortcut from "./shortcut.vue";

    export default {
        components: {emviCmdIcon, emviCmdShortcut},
        props: {
            icon: {default: ""},
            color: {default: ""},
            view: {default: ""},
            index: {default: 0},
            showShortcut: {default: false},
            disabled: {default: false},
            enter: {default: false},
            tab: {default: false}
        },
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.index === this.row;
            },
            submenuClass() {
                return {
                    entry: true,
                    active: this.active && !this.isTouch,
                    disabled: this.disabled
                };
            },
            colorClass() {
                let c = {"button cursor-pointer": true};
                c[`${this.color}-100`] = this.color;
                c[`bg-${this.color}-10`] = this.color;
                return c;
            }
        },
        watch: {
            row() {
                if(this.active) {
                    this.focusCmdInput();
                    scrollIntoViewArea(this.$refs.submenu, document.getElementById("cmd-results"));
                }
            },
            enter(enter) {
                if(enter && this.active) {
                    this.open();
                }
            },
            // Usually we use tab/up/down outside of this component to change the active row,
            // but it's sometimes useful when this component is used inside a form (with input, select, ... elements).
            tab(tab) {
                if(tab && this.active) {
                    if(tab.shiftKey) {
                        this.$emit("previous");
                    }
                    else {
                        this.$emit("next");
                    }
                }
            }
        },
        methods: {
            open() {
                if(this.disabled) {
                    return;
                }

                if(!this.view) {
                    this.$emit("enter");
                }
                else {
                    this.$store.dispatch("pushColumn", this.view);
                }

                this.focusCmdInput();
            }
        }
    }
</script>
