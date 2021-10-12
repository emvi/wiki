<template>
    <div class="avatar size-48 cursor-pointer" v-on:click="$emit('click')">
        <img :src="picture" alt="" v-if="picture" />
        <div class="initials pink-100 bg-pink-10" v-if="!picture">{{user | initials}}</div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        computed: {
            ...mapGetters(["user"]),
            picture() {
                return this.user.picture;
            }
        }
    }
</script>
