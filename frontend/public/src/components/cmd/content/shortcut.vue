<template>
    <span class="shortcut cursor-pointer" v-on:click.stop="$emit('click')">
        <span v-if="!isTouch">
            <span class="key">{{shortcut}}</span>
            <span>
                <slot></slot>
            </span>
        </span>
        <i :class="iconClass" v-if="isTouch && icon"></i>
    </span>
</template>

<script>
    export default {
        props: ["shortcut", "icon", "rotate"],
        computed: {
            iconClass() {
                let c = {
                    "icon action": true,
                    "icon-rotate-180": this.rotate
                };
                c[`icon-${this.icon}`] = true;
                return c;
            }
        }
    }
</script>
