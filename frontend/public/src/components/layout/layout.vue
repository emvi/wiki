<template>
    <div>
        <emvi-header v-if="!isMobile"></emvi-header>
        <emvi-mobile-header v-if="isMobile"></emvi-mobile-header>
        <emvi-cmd></emvi-cmd>
        <emvi-navigation :hide="hideNavigation" v-if="!isMobile"></emvi-navigation>
        <emvi-actions :hide="hideNavigation" v-if="!isMobile"></emvi-actions>
        <main>
            <slot></slot>
        </main>
        <emvi-footer v-if="!isMobile"></emvi-footer>
        <emvi-toast></emvi-toast>
    </div>
</template>

<script>
    import emviHeader from "./header.vue";
    import emviFooter from "./footer.vue";
    import emviCmd from "../cmd/cmd.vue";
    import emviToast from "../toast/toast.vue";
    import emviMobileHeader from "../mobile/header.vue";
    import emviNavigation from "./navigation.vue";
    import emviActions from "./actions.vue";

    export default {
        components: {
            emviHeader,
            emviFooter,
            emviCmd,
            emviToast,
            emviMobileHeader,
            emviNavigation,
            emviActions
        },
        props: {
            disableEvents: {default: false},
            hideNavigation: {default: false},
            hideActions: {default: false}
        },
        data() {
            return {
                keydownHandler: null
            };
        },
        mounted() {
            if(!this.disableEvents) {
                this.bindKeys();
            }
        },
        beforeDestroy() {
            if(!this.disableEvents) {
                this.unbindKeys();
            }
        },
        methods: {
            bindKeys() {
                this.keydownHandler = e => {
                    let prevent = true;

                    switch(e.code) {
                        case "ArrowUp":
                            this.$emit("up", e);
                            break;
                        case "ArrowDown":
                            this.$emit("down", e);
                            break;
                        case "ArrowLeft":
                            this.$emit("left", e);
                            prevent = false;
                            break;
                        case "ArrowRight":
                            this.$emit("right", e);
                            prevent = false;
                            break;
                        case "Tab":
                            this.$emit("tab", e);
                            break;
                        case "Enter":
                            this.$emit("enter", e);
                            break;
                        case "Escape":
                            this.$emit("esc", e);
                            break;
                        default:
                            prevent = false;
                    }

                    if(prevent) {
                        e.preventDefault();
                        e.stopPropagation();
                    }
                };
                window.addEventListener("keydown", this.keydownHandler);
            },
            unbindKeys() {
                window.removeEventListener("keydown", this.keydownHandler);
            }
        }
    }
</script>
