<template>
    <div :class="cardClass" :scroll-area="scrollArea" v-on:click="click" v-on:mouseenter="$emit('mouseenter')" v-on:mouseleave="$emit('mouseleave')" ref="card">
        <div class="card-content">
            <emvi-avatar :entity="entity" :icon="icon" size="40" v-on:click="clickAvatar"></emvi-avatar>
            <slot></slot>
        </div>
        <slot name="after"></slot>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {scrollIntoView, scrollIntoViewArea} from "../../util";
    import emviShortcut from "../cmd/content/shortcut.vue";
    import emviAvatar from "../cmd/content/avatar.vue";

    export default {
        components: {emviShortcut, emviAvatar},
        props: {
            icon: "",
            entity: {},
            active: {default: false},
            disabled: {default: false},
            scrollArea: ""
        },
        computed: {
            ...mapGetters(["selection"]),
            cardClass() {
                let c = {};
                c["card cursor-pointer"] = true;
                c["focus-ring"] = this.active && !this.isTouch;
                c["disabled"] = this.disabled;
                return c;
            },
            iconClass() {
                let c = {icon: true};
                c[`icon-${this.icon}`] = true;
                return c;
            }
        },
        watch: {
            selection() {
                if(this.active) {
                    if(this.scrollArea) {
                        scrollIntoViewArea(this.$refs.card, this.scrollArea);
                    }
                    else {
                        scrollIntoView(this.$refs.card);
                    }
                }
            }
        },
        methods: {
            click(e) {
                this.$emit("click", e);
            },
            clickAvatar(e) {
                this.click(e);
                this.$emit("iconclick");
            }
        }
    }
</script>
