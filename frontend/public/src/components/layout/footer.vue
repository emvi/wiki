<template>
    <footer class="no-select">
        <i class="icon icon-help size-48 cursor-pointer" v-on:click.stop.prevent="open"></i>
    </footer>
</template>

<script>
    export default {
        methods: {
            open() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("help_cmd"));
                    this.$store.dispatch("pushColumn", "help");
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "help_cmd": "/help"
        },
        "de": {
            "help_cmd": "/hilfe"
        }
    }
</i18n>
