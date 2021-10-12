<template>
    <transition name="expand" v-on:enter="enter" v-on:after-enter="afterEnter" v-on:leave="leave">
        <slot></slot>
    </transition>
</template>

<script>
    export default {
        methods: {
            afterEnter(element) {
                element.style.height = 'auto';
            },
            enter(element) {
                element.style.width = getComputedStyle(element).width;
                element.style.position = 'absolute';
                element.style.visibility = 'hidden';
                element.style.height = 'auto';

                const height = getComputedStyle(element).height;
                element.style.width = null;
                element.style.position = null;
                element.style.visibility = null;
                element.style.height = 0;

                // force repaint to make sure the animation is triggered correctly
                getComputedStyle(element).height;

                setTimeout(() => {
                    element.style.height = height;
                });
            },
            leave(element) {
                element.style.height = getComputedStyle(element).height;

                // force repaint to make sure the animation is triggered correctly
                getComputedStyle(element).height;

                setTimeout(() => {
                    element.style.height = 0;
                });
            }
        }
    }
</script>

<style lang="scss" scoped>
    * {
        will-change: height;
        transform: translateZ(0);
        backface-visibility: hidden;
        perspective: 1000px;
    }

    .expand-enter-active, .expand-leave-active {
        scrollbar-color: transparent transparent;

        &::-webkit-scrollbar-thumb {
            background-color: transparent;
        }
    }

    span.expand-enter-active, span.expand-leave-active {
        display: block;
    }
</style>
