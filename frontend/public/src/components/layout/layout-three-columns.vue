<template>
    <div class="layout narrow">
        <div class="topbar right no-select" v-if="!isMobile">
            <slot name="top"></slot>
        </div>
        <div class="sidebar left no-select" v-show="$mq > 1440">
            <div class="sidebar-content">
                <slot name="left"></slot>
            </div>
        </div>
        <slot></slot>
        <div class="sidebar right no-select" v-show="$mq > 1440">
            <div class="sidebar-content">
                <slot name="right"></slot>
            </div>
        </div>
    </div>
</template>

<script>
    export default {}
</script>
