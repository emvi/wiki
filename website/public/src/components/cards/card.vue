<template>
    <div v-bind:class="{'card--collapse': (!collapsible || collapsed) && !childColor, 'card--expand': (collapsible && !collapsed) && !childColor, 'card--unread': unread, 'card--blue': color === 'blue' && !childColor, 'card--purple': color === 'purple' && !childColor, 'card--pink': color === 'pink' && !childColor, 'card--green': color === 'green' && !childColor}">
        <div v-bind:class="{'card--content': !childColor, 'card--content--blue': childColor === 'blue', 'card--content--pink': childColor === 'pink'}">
            <div v-bind:class="{'card--content--left': true, 'cursor-default': cursorDefault}" v-on:click="toggleCollapsed">
                <slot name="icon"></slot>
            </div>
            <div class="card--content--center">
                <div class="card--content--center--info" v-on:click="toggleCollapsed" v-if="showInfo">
                    <slot name="info"></slot>
                </div>
                <span v-bind:class="{'card--content--center--item': true, 'cursor-default': cursorDefault}" v-on:click="toggleCollapsed" v-if="showTitle">
                    <slot name="title"></slot>
                </span>
                <slot></slot>
            </div>
            <div class="card--content--right" v-if="showRight || collapsible">
                <slot name="right"></slot>
                <i v-bind:class="{'icon card--content--icon icon-expand': true, 'icon-rotate-180': !collapsed}" v-if="collapsible" v-on:click="toggleCollapsed"></i>
            </div>
        </div>
        <expand>
            <div v-bind:class="{'card--content--list': maxheight}" v-show="!collapsed">
                <slot name="children"></slot>
            </div>
        </expand>
    </div>
</template>

<script>
    import expand from "../expand.vue";

    export default {
        components: {
            expand
        },
        props: {
            collapsible: {default: false},
            color: {default: "blue"},
            childColor: {default: null},
            maxheight: {default: true},
            cursorDefault: {default: false},
            unread: {default: false},
            showTitle: {default: false},
            showInfo: {default: false},
            showRight: {default: false},
        },
        data() {
            return {
                collapsed: true,
            };
        },
        methods: {
            toggleCollapsed() {
                if(this.collapsible) {
                    this.collapsed = !this.collapsed;

                    if (!this.collapsed) {
                        this.$emit("expand");
                    }
                }
            },
            collapse() {
                this.collapsed = true;
            }
        }
    }
</script>
