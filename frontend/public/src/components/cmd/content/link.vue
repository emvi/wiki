<template>
    <div :class="linkClass" v-on:click.stop.prevent="click">
        <emvi-cmd-icon icon="link"></emvi-cmd-icon>
        <div class="item">
            <p v-show="link">{{link}}</p>
            <p v-show="!link">{{$t("placeholder")}}</p>
        </div>
        <emvi-cmd-shortcut shortcut="Enter" v-show="link && active" v-on:click="click">
            {{$t("copy")}}
        </emvi-cmd-shortcut>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import emviCmdIcon from "./icon.vue";
    import emviCmdShortcut from "./shortcut.vue";

    export default {
        components: {emviCmdIcon, emviCmdShortcut},
        props: {
            index: {default: 0},
            link: {default: ""}
        },
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.row === this.index;
            },
            linkClass() {
                return {
                    "entry cursor-pointer": true,
                    active: this.active && !this.isTouch
                };
            }
        },
        methods: {
            click() {
                this.$emit("click");
                this.$store.dispatch("selectRow", this.index);
                this.focusCmdInput();
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder": "You need to generate a link first.",
            "copy": "Copy to clipboard"
        },
        "de": {
            "placeholder": "Du musst zunächst einen Link generieren.",
            "copy": "In Zwischenablage kopieren"
        }
    }
</i18n>
