<template>
    <div :class="colorClass" v-on:click.stop.prevent="click">
        <div :class="{'entry': true, 'active': showActive}">
            <emvi-cmd-icon :icon="icon"></emvi-cmd-icon>
            {{label}}
        </div>
    </div>
</template>

<script>
    import {FormMixin} from "../mixin";
    import emviCmdIcon from "../../content/icon.vue";

    export default {
        mixins: [FormMixin],
        components: {emviCmdIcon},
        props: {
            icon: {default: ""},
            label: {default: ""},
            color: {default: ""}
        },
        computed: {
            colorClass() {
                let c = {"button cursor-pointer": true};
                c[`${this.color}-100`] = true;
                c[`bg-${this.color}-10`] = true;
                return c;
            }
        },
        methods: {
            click() {
                this.$emit("click");
                this.focus();
            }
        }
    }
</script>
