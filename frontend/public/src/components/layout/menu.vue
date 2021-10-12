<template>
    <div class="menu">
        <i class="icon icon-add size-48 cursor-pointer" v-on:click.stop.prevent="create" v-if="member.show_create_button && !isReadOnly"></i>
        <emvi-notifications></emvi-notifications>
        <emvi-avatar v-on:click="openProfile"></emvi-avatar>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import emviNotifications from "./notifications.vue";
    import emviAvatar from "./avatar.vue";

    export default {
        components: {emviNotifications, emviAvatar},
        computed: {
            ...mapGetters(["member"])
        },
        methods: {
            create() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("create_cmd"));
                    this.$store.dispatch("pushColumn", "create");
                });
            },
            openProfile() {
                window.open(`${EMVI_WIKI_WEBSITE_HOST}/account`, "_blank");
            }
        }
    }
</script>

<i18n>
{
    "en": {
        "create_cmd": "/new"
    },
    "de": {
        "create_cmd": "/neu"
    }
}
</i18n>
