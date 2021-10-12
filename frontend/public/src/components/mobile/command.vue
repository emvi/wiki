<template>
    <i class="icon icon-flash size-48" v-on:click.stop="open"></i>
</template>

<script>
    export default {
        methods: {
            open() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", "/");
                });
            }
        }
    }
</script>
