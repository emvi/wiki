<template>
    <i :class="iconClass" v-on:click="watch"></i>
</template>

<script>
    import {toggleObserve} from "../cmd/command_funcs";

    export default {
        props: ["observed", "type"],
        computed: {
            iconClass() {
                let c = {"icon icon-notification cursor-pointer": true};
                c[`grey-50`] = !this.observed;
                c[`yellow-100`] = this.observed;
                return c;
            }
        },
        methods: {
            watch() {
                toggleObserve(this, this.type);
            }
        }
    }
</script>
