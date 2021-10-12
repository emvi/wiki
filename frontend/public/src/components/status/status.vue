<template>
    <small class="status no-select" :title="title">
        <i :class="iconClass"></i>
        <span class="grey-60-to-grey-40">
            <slot></slot>
        </span>
    </small>
</template>

<script>
    export default {
        props: ["icon", "title"],
        computed: {
            iconClass() {
                let c = {"icon size-32 grey-50": true};
                c[`icon-${this.icon}`] = true;
                return c;
            }
        }
    }
</script>
