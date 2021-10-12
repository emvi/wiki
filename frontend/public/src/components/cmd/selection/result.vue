<template>
    <div :class="{'cursor-pointer': cursorPointer}" v-on:click.stop="click" v-on:mouseenter="$emit('mouseenter')" v-on:mouseleave="$emit('mouseleave')" ref="result">
        <div :class="resultClass">
            <emvi-cmd-avatar :entity="entity" :icon="icon"></emvi-cmd-avatar>
            <slot></slot>
        </div>
        <div class="preview">
            <div class="preview-content" id="cmd-selection-result-details" v-if="hasDetails && details">
                <slot name="details"></slot>
            </div>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {scrollIntoViewArea, scrollToTopArea} from "../../../util";
    import emviCmdAvatar from "../content/avatar.vue";

    export default {
        components: {emviCmdAvatar},
        props: {
            icon: {default: ""},
            cursorPointer: {default: false},
            entity: {default: () => {return {}}},
            index: {default: 0, required: true},
            disabled: {default: false},
            details: {default: false}
        },
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.index === this.row;
            },
            hasDetails() {
                return !!this.$slots.details;
            },
            resultClass() {
                return {
                    entry: true,
                    active: this.active && !this.isTouch,
                    disabled: this.disabled,
                    double: this.isMobile
                };
            }
        },
        watch: {
            row() {
                if(this.active) {
                    scrollIntoViewArea(this.$refs.result, document.getElementById("cmd-results"));
                }
            },
            details(details) {
                this.$nextTick(() => {
                    if(details) {
                        scrollToTopArea(this.$refs.result, document.getElementById("cmd-results"));
                    }
                    else {
                        scrollIntoViewArea(this.$refs.result, document.getElementById("cmd-results"));
                    }
                });
            }
        },
        methods: {
            click() {
                this.$emit("click");
                let search = document.getElementById("cmd-search-input");

                if(search) {
                    search.focus();
                }
                else {
                    this.focusCmdInput();
                }
            }
        }
    }
</script>
