<template>
    <emvi-card :icon="icon" :entity="entity" :active="active" :disabled="disabled" :scroll-area="scrollArea" v-on:click="click" v-on:mouseenter="mouseEnter" v-on:mouseleave="mouseLeave" ref="card">
        <template>
            <div class="item">
                <p>
                    <slot></slot>
                </p>
                <small>
                    <slot name="info"></slot>
                </small>
            </div>
            <emvi-shortcut shortcut="Tab" v-show="showShortcut" v-on:click="$emit('preview')">
                {{$t("preview")}}
            </emvi-shortcut>
            <i :class="{'icon icon-chevron action cursor-pointer': true, 'icon-rotate-180': preview}" v-show="showPreviewArrow" v-on:click.stop="$emit('preview')"></i>
        </template>
        <template slot="after">
            <div class="preview" v-show="showPreview">
                <div class="preview-content" ref="preview">
                    <slot name="preview"></slot>
                </div>
            </div>
        </template>
    </emvi-card>
</template>

<script>
    import {mapGetters} from "vuex";
    import emviCard from "./card.vue";
    import emviShortcut from "../cmd/content/shortcut.vue";
    import {scrollArea} from "../../util";

    export default {
        components: {emviCard, emviShortcut},
        props: {
            icon: "",
            entity: {},
            active: {default: false},
            disabled: {default: false},
            preview: {default: false},
            scrollArea: "",
            up: {default: false},
            down: {default: false}
        },
        data() {
            return {
                hover: false
            };
        },
        computed: {
            ...mapGetters(["selection"]),
            showShortcut() {
                return (this.active || this.hover) && !this.isTouch;
            },
            showPreviewArrow() {
                return this.active || this.hover || this.isTouch;
            },
            showPreview() {
                return this.active && this.preview;
            }
        },
        watch: {
            up(up) {
                if(up && this.active && this.preview) {
                    scrollArea(this.$refs.preview, -1);
                }
            },
            down(down) {
                if(down && this.active && this.preview) {
                    scrollArea(this.$refs.preview, 1);
                }
            }
        },
        methods: {
            click(e) {
                this.$emit("click", e);
            },
            mouseEnter() {
                this.hover = true;
            },
            mouseLeave() {
                this.hover = false;
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
