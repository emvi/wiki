<template>
    <div :class="{'entry cursor-pointer': true, 'active': showActive}" v-on:click.stop.prevent="click" ref="entry">
        <div class="avatar size-32">
            <i class="icon icon-seen"></i>
        </div>
        <span class="item">
            <p>{{$t("text")}}</p>
        </span>        
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {scrollIntoViewArea} from "../../../../util";

    export default {
        props: ["index"],
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.index === this.row;
            },
            showActive() {
                return this.active && !this.isTouch;
            }
        },
        watch: {
            row() {
                if(this.active) {
                    scrollIntoViewArea(this.$refs.entry, document.getElementById("cmd-results"));
                }
            }
        },
        methods: {
            click() {
                this.$emit('click');
                this.focusCmdInput();
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text": "Mark all as read"
        },
        "de": {
            "text": "Alle als gelesen markieren"
        }
    }
</i18n>
