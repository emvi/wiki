<template>
    <div :class="avatarClass" v-on:click.stop.prevent="$emit('click')">
        <i :class="iconClass" v-if="!isUser"></i>
        <img :src="picture" alt="" v-if="isUser && picture" />
        <div class="initials pink-100 bg-pink-10" v-if="isUser && !picture">{{entity | initials}}</div>
    </div>
</template>

<script>
    export default {
        props: {
            size: {default: "32"},
            icon: {default: ""},
            color: "",
            entity: {default: () => {return {}}}
        },
        computed: {
            avatarClass() {
                let c = {"avatar cursor-pointer": true};
                c[`size-${this.size}`] = true;

                if(this.color) {
                    c[`${this.color}-100 bg-${this.color}-10`] = true;
                }

                return c;
            },
            iconClass() {
                let c = {icon: true};
                c[`icon-${this.icon}`] = true;
                return c;
            },
            isUser() {
                return this.entity.type === "user";
            },
            picture() {
                return this.entity.picture || "";
            }
        }
    }
</script>
