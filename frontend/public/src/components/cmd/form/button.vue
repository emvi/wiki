<template>
    <div :class="colorClass" v-on:click.stop.prevent="click" ref="entry">
        <div :class="{'entry': true, 'active': showActive, 'disabled': disabled}" tabindex="0"
             ref="element"
             v-on:keydown.tab.stop.prevent="tab"
             v-on:keydown.enter.stop.prevent="$emit('enter')"
             v-on:keydown.esc.stop.prevent="$emit('esc')"
             v-on:keydown.up.stop="up"
             v-on:keydown.down.stop="down">
            <div class="avatar size-32">
                <i :class="iconClass"></i>
            </div>
            <div class="item">
                <p>{{label}}</p>
            </div>
        </div>
    </div>
</template>

<script>
    import {FormMixin} from "./mixin";

    export default {
        mixins: [FormMixin],
        props: {
            color: {default: "green"}
        },
        computed: {
            colorClass() {
                let c = {"button cursor-pointer": true};
                c[`${this.color}-100`] = true;
                c[`bg-${this.color}-10`] = true;
                return c;
            },
            iconClass() {
                let c = {icon: true};
                c[`icon-${this.icon}`] = true;
                return c;
            }
        },
        methods: {
            up() {
                this.$emit("previous");
            },
            down() {
                this.$emit("next");
            },
            click() {
                this.$emit("enter");
                this.setSelected();
                this.focusCmdInput();
            }
        }
    }
</script>
