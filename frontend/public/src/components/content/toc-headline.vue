<template>
    <div class="toc">
        <span :class="{'level-1 cursor-pointer': true, 'active': headline.active}" v-on:click="jump(headline)" v-if="headline.tag === 'H2'">
            <i>{{headline.number}}</i>
            <p>{{headline.text}}</p>
        </span>
        <toc v-for="h in headline.children"
             :key="h.id"
             :headline="h"
             v-on:jump="jump"></toc>
        <small :class="{'cursor-pointer': true, 'active': headline.active, 'indent': headline.tag === 'H4'}" v-on:click="jump(headline)" v-if="headline.tag === 'H3' || headline.tag === 'H4'">
            {{headline.text}}
        </small>
    </div>
</template>

<script>
    export default {
        name: "toc",
        props: ["headline"],
        methods: {
            jump(headline) {
                this.$emit("jump", headline);
            }
        }
    }
</script>
