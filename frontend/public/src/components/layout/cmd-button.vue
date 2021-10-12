<template>
    <div :class="{'cmd-button cursor-pointer no-select': true, 'visible': show}" v-on:click.stop="openCmd" v-if="!cmdOpen" ref="button">
        <i class="icon icon-circle-12 size-40"></i>
        <small>{{$t("cmd")}}</small>
        <emvi-cmd-shortcut :shortcut="$t('shift')" v-on:click="openCmd" v-if="!isTouch">
            <span>&#43;</span>
            <span class="key">{{$t("space")}}</span>
        </emvi-cmd-shortcut>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import emviCmdShortcut from "../cmd/content/shortcut.vue";
    import {debounce} from "../../util";

    const hideTimeout = 500;

    export default {
        components: {emviCmdShortcut},
        data() {
            return {
                mouseHandler: null,
                hideDebounced: null,
                show: true
            };
        },
        computed: {
            ...mapGetters(["cmdOpen"])
        },
        mounted() {
            this.hideDebounced = debounce(() => {
                this.show = false;
            }, hideTimeout);
            this.hideDebounced();
            this.bindMouse();
        },
        beforeDestroy() {
            this.unbindMouse();
        },
        methods: {
            bindMouse() {
                this.mouseHandler = document.addEventListener("mousemove", () => {
                    this.show = true;
                    this.hideDebounced();
                });
            },
            unbindMouse() {
                document.removeEventListener("mousemove", this.mouseHandler);
            },
            openCmd() {
                this.$store.commit("setCmdOpen", true);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "cmd": "Type a command or search...",
            "shift": "Shift",
            "space": "Space"
        },
        "de": {
            "cmd": "Tippe einen Befehl oder suche...",
            "shift": "Shift",
            "space": "Leer"
        }
    }
</i18n>
