<template>
    <div class="profile-header no-select">
        <emvi-cmd-avatar size="96" :entity="entity"></emvi-cmd-avatar>
        <h1>
            <slot></slot>
        </h1>
    </div>
</template>

<script>
    import emviCmdAvatar from "../content/avatar.vue";

    export default {
        components: {emviCmdAvatar},
        props: ["entity"]
    }
</script>
