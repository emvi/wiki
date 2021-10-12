<template>
    <div>
        <emvi-cmd-notification v-for="(notification, index) in notifications"
            :key="notification.id"
            :notification="notification"
            :read="notification.read"
            :index="index"
            :tab="tab"
            :enter="enter"></emvi-cmd-notification>
        <emvi-cmd-empty v-show="!notifications.length"></emvi-cmd-empty>
        <emvi-cmd-notification-all-read :index="notifications.length" v-on:click="markAllRead"></emvi-cmd-notification-all-read>
        <emvi-cmd-all-results :index="notifications.length+1" v-on:click="showAll"></emvi-cmd-all-results>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../util";
    import {FeedService} from "../../../service";
    import emviCmdNotification from "./notifications/notification.vue";
    import emviCmdAllResults from "./notifications/all-results.vue";
    import emviCmdEmpty from "../content/empty.vue";
    import emviCmdNotificationAllRead from "./notifications/notification-all-read.vue";

    export default {
        components: {
            emviCmdNotification,
            emviCmdAllResults,
            emviCmdEmpty,
            emviCmdNotificationAllRead
        },
        props: ["up", "down", "enter", "tab", "esc"],
        computed: {
            ...mapGetters(["row", "notifications"]),
        },
        watch: {
            row(row) {
                updateSelectedRow(row, this.notifications.length+2, this.$store);
            },
            up(up) {
                if(up) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down) {
                    this.$store.dispatch("selectNextRow");
                }
            },
            enter(enter) {
                if(enter) {
                    if(this.row === this.notifications.length) {
                        this.markAllRead();
                    }
                    else if(this.row === this.notifications.length+1) {
                        this.showAll();
                    }
                }
            },
            esc(esc) {
                if(esc) {
                    if(this.$route.name === "activities") {
                        this.$store.dispatch("setMeta", {key: "updateList", value: true});
                    }

                    this.$store.dispatch("popColumn");
                }
            }
        },
        methods: {
            markAllRead() {
                this.resetError();

                if(!this.notifications.length) {
                    return;
                }

                FeedService.toggleNotificationRead()
                    .then(() => {
                        this.$store.dispatch("markNotificationsRead");
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            showAll() {
                if(this.$route.name !== "notifications") {
                    this.$router.push("/notifications");
                }

                this.$store.dispatch("closeCmd");
            }
        }
    }
</script>
