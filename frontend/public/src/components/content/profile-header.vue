<template>
    <div class="profile-header no-select">
        <div class="avatar size-96">
            <img :src="picture" alt="" v-if="picture" />
            <span v-if="!picture">{{entity | initials}}</span>
        </div>
        <h1>
            <slot></slot>
        </h1>
    </div>
</template>

<script>
    export default {
        props: ["entity"],
        computed: {
            picture() {
                return this.entity.picture;
            }
        }
    }
</script>
