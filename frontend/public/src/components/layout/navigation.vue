<template>
    <div :class="{'drawer left': true, 'expand': !hidden, 'no-transition': noTransition}" v-if="member.show_navigation" ref="navigation">
        <i class="toggle icon icon-compass cursor-pointer" v-on:click="toggle"></i>
        <div class="drawer-content focus-ring">
            <router-link v-for="nav in navigation" :key="nav.pageName" :to="nav.path" v-if="!nav.requires_admin || isAdmin">
                <i :class="getIconClass(nav.icon)"></i>
                {{capitalize(nav.page[locale][0])}}
            </router-link>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {capitalize} from "../../util";
    import {navigation} from "../cmd/navigation";

    export default {
        props: {
            hide: {default: false}
        },
        data() {
            return {
                mouseHandler: null,
                navigation,
                noTransition: true,
                hidden: true
            };
        },
        computed: {
            ...mapGetters(["member"])
        },
        mounted() {
            this.hidden = this.hide || this.$mq <= 1440;

            if(this.hidden) {
                this.bindMouse();
            }
        },
        beforeDestroy() {
            if(this.mouseHandler) {
                this.unbindMouse();
            }
        },
        methods: {
            bindMouse() {
                this.mouseHandler = document.addEventListener("mouseup", e => {
                    if(this.$refs.navigation && e.target !== this.$refs.navigation && !this.$refs.navigation.contains(e.target)) {
                        this.hidden = true;
                    }
                });
            },
            unbindMouse() {
                document.removeEventListener("mouseup", this.mouseHandler);
            },
            capitalize(name) {
                return capitalize(name);
            },
            getIconClass(icon) {
                let c = {icon: true};
                c[`icon-${icon} size-40`] = true;
                return c;
            },
            toggle() {
                this.noTransition = false;
                this.hidden = !this.hidden;
            }
        }
    }
</script>
