<template>
    <div :class="userClass" :title="title">
        <emvi-avatar size="40" :entity="author"></emvi-avatar>
    </div>
</template>

<script>
    import emviAvatar from "../cmd/content/avatar.vue";

    export default {
        components: {emviAvatar},
        props: {
            author: {},
        },
        computed: {
            userClass() {
                let c = {user: true};
                c[`${this.author.color}`] = true;
                return c;
            },
            title() {
                return `${this.author.firstname} ${this.author.lastname}`;
            }
        }
    }
</script>
