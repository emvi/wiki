<template>
    <div class="avatar size-32">
        <i :class="iconClass"></i>
    </div>
</template>

<script>
    export default {
        props: ["icon"],
        computed: {
            iconClass() {
                let c = {icon: true};
                c[`icon-${this.icon}`] = true;
                return c;
            }
        }
    }
</script>
