<template>
    <div class="entry">
        <emvi-cmd-avatar :entity="entity" :icon="icon"></emvi-cmd-avatar>
        <slot></slot>
    </div>
</template>

<script>
    import emviCmdAvatar from "../content/avatar.vue";

    export default {
        components: {emviCmdAvatar},
        props: {
            entity: {},
            icon: ""
        }
    }
</script>
