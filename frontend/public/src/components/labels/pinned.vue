<template>
    <i :class="iconClass" v-on:click="pin"></i>
</template>

<script>
    import {togglePinned} from "../cmd/command_funcs";

    export default {
        props: ["pinned", "type"],
        computed: {
            iconClass() {
                let c = {"icon icon-pin cursor-pointer": true};
                c[`grey-50`] = !this.pinned;
                c[`brown-100`] = this.pinned;
                return c;
            }
        },
        methods: {
            pin() {
                togglePinned(this, this.type);
            }
        }
    }
</script>
