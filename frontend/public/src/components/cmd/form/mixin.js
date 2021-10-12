import {mapGetters} from "vuex";
import {scrollIntoViewArea} from "../../../util";

export const FormMixin = {
    props: {
        container: {default: "cmd-results"},
        index: {default: 0},
        value: {default: ""},
        disabled: {default: false},
        hint: {default: ""},
        error: {default: ""},
        required: {default: false},
        optional: {default: false},
        label: "",
        icon: ""
    },
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
        active() {
            this.focus();
        }
    },
    mounted() {
        this.focus();
    },
    methods: {
        focus() {
            if(this.active) {
                this.focusElement();
                scrollIntoViewArea(this.$refs.entry, document.getElementById(this.container));
            }
        },
        focusElement() {
            this.$refs.element.focus();
        },
        setSelected() {
            this.$store.dispatch("selectRow", this.index);
        },
        tab(e) {
            if(e.shiftKey) {
                this.$emit("previous");
            }
            else {
                this.$emit("next");
            }
        },
        up() {
            // do nothing, so that this behaviour can be overwritten
        },
        down() {
            // do nothing, so that this behaviour can be overwritten
        }
    }
};
