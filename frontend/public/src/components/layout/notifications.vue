<template>
    <i class="icon icon-notification size-48 cursor-pointer" v-on:click.stop.prevent="open">
        <span class="counter" v-if="notificationCount">{{notificationCount}}</span>
    </i>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        data() {
            return {
                notificationInterval: null
            };
        },
        computed: {
            ...mapGetters(["view", "notificationCount"])
        },
        mounted() {
            this.loadNotifications();
        },
        beforeDestroy() {
            clearInterval(this.notificationInterval);
        },
        methods: {
            loadNotifications() {
                this.$store.dispatch("loadNotifications");

                this.notificationInterval = setInterval(() => {
                    if(this.view !== "notifications") {
                        this.$store.dispatch("loadNotifications");
                    }
                }, 5000);
            },
            open() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("notifications_cmd"));
                    this.$store.dispatch("pushColumn", "notifications");
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "notifications_cmd": "/notifications"
        },
        "de": {
            "notifications_cmd": "/benachrichtigungen"
        }
    }
</i18n>
