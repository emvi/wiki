<template>
    <div class="card-content" v-on:click.stop.prevent="$emit('click')">
        <emvi-avatar size="40" :entity="entity" :icon="icon" v-on:click="$emit('click')"></emvi-avatar>
        <slot></slot>
    </div>
</template>

<script>
    import emviAvatar from "../cmd/content/avatar.vue";

    export default {
        components: {emviAvatar},
        props: {
            entity: {},
            icon: ""
        }
    }
</script>
