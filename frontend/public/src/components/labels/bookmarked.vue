<template>
    <i :class="iconClass" v-on:click="bookmark"></i>
</template>

<script>
    import {toggleBookmark} from "../cmd/command_funcs";

    export default {
        props: ["bookmarked", "type"],
        computed: {
            iconClass() {
                let c = {"icon icon-bookmark cursor-pointer": true};
                c[`grey-50`] = !this.bookmarked;
                c[`bluegreen-100`] = this.bookmarked;
                return c;
            }
        },
        methods: {
            bookmark() {
                toggleBookmark(this, this.type);
            }
        }
    }
</script>
