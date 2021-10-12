<template>
    <header>
        <emvi-orga></emvi-orga>
        <emvi-cmd-button></emvi-cmd-button>
        <emvi-menu></emvi-menu>
    </header>
</template>

<script>
    import emviOrga from "./orga.vue";
    import emviMenu from "./menu.vue";
    import emviCmdButton from "./cmd-button.vue";

    export default {
        components: {emviOrga, emviMenu, emviCmdButton}
    }
</script>
