<template>
    <div class="enlarge" v-on:click="close">
        <div class="article-content" >
            <div ref="content">
                <slot></slot>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: ["content"],
        methods: {
            close(e) {
                if(!this.$refs.content.contains(e.target)) {
                    this.$emit("close");
                }
            }
        }
    }
</script>
