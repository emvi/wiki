<template>
    <i class="icon icon-home size-48" v-on:click.stop="toStart"></i>
</template>

<script>
    export default {
        methods: {
            toStart() {
                if(this.$route.name !== "start") {
                    this.$store.dispatch("resetCmd");
                    this.$router.push("/");
                }
            },
        }
    }
</script>
