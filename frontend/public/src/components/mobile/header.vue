<template>
    <header>
        <emvi-mobile-start></emvi-mobile-start>
        <emvi-mobile-search></emvi-mobile-search>
        <emvi-mobile-command></emvi-mobile-command>
        <emvi-mobile-navigation></emvi-mobile-navigation>
        <emvi-notifications></emvi-notifications>
    </header>
</template>

<script>
    import emviMobileStart from "../mobile/start.vue";
    import emviMobileSearch from "../mobile/search.vue";
    import emviMobileCommand from "../mobile/command.vue";
    import emviMobileNavigation from "../mobile/navigation.vue";
    import emviNotifications from "../layout/notifications.vue";

    export default {
        components: {
            emviMobileStart,
            emviMobileSearch,
            emviMobileCommand,
            emviMobileNavigation,
            emviNotifications
        }
    }
</script>
