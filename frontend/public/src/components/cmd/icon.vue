<template>
    <i :class="iconClass" v-on:click="$emit('click')" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false"></i>
</template>

<script>
    import {mapGetters} from "vuex";
    import {icons} from "./icons";

    export default {
        data() {
            return {
                iconClass: {"icon icon-circle-12 cursor-pointer": true},
                hover: false
            };
        },
        computed: {
            ...mapGetters(["view", "columns"])
        },
        watch: {
            view(view) {
                this.setIcon(view);
            },
            hover(hover) {
                if(hover && this.columns > 0) {
                    this.setIconBack();
                }
                else {
                    this.setIcon(this.view);
                }
            }
        },
        beforeMount() {
            this.setIcon(this.view);
        },
        methods: {
            setIcon(view) {
                if(this.isTouch && this.columns > 1) {
                    this.setIconBack();
                    return;
                }

                if(!view && this.columns === 0) {
                    view = "default";
                }

                let iconClass = icons.get(view);
                let icon = {"icon cursor-pointer": true};

                if(iconClass) {
                    icon[iconClass] = true;
                }
                else {
                    icon["icon-circle-12"] = true;
                }

                this.iconClass = icon;
            },
            setIconBack() {
                this.iconClass = {"icon icon-back cursor-pointer": true};
            }
        }
    }
</script>
