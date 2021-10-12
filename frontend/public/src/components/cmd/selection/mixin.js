import {mapGetters} from "vuex";
import {updateSelection} from "../util";

export const PropsMixin = {
    props: {
        entity: {},
        index: {default: 0, required: true},
        disabled: {default: false},
        enter: {default: false},
        tab: {default: false},
        del: {default: false},
        esc: {default: false},
        up: {default: false},
        down: {default: false}
    },
    data() {
        return {
            hover: false,
            details: false
        };
    },
    computed: {
        ...mapGetters(["row"]),
        active() {
            return this.row === this.index || this.hover || this.isTouch;
        }
    },
    watch: {
        row() {
            if(!this.active && this.details) {
                this.details = false;
            }
        },
        details(details) {
            if(!details && this.selection) {
                this.selection = 0;
            }
        }
    },
    methods: {
        emitDetails(e) {
            this.$emit("details", e);
        },
        emitRemove(e) {
            this.$emit("remove", e);
        },
        toggleDetails(showDetails) {
            this.details = showDetails;
            this.$emit("details", showDetails);
            this.$store.dispatch("selectRow", this.index);
            this.focusCmdInput();
        }
    }
};

export const SelectionMixin = {
    mixins: [PropsMixin],
    data() {
        return {
            selection: 0
        };
    },
    watch: {
        up(up) {
            if(up && this.active && this.details) {
                this.previousRow();
            }
        },
        down(down) {
            if(down && this.active && this.details) {
                this.nextRow();
            }
        }
    },
    methods: {
        nextRow() {
            this.selection++;
            this.selection = updateSelection(this.selection, this.maxSelectionIndex || 0);
        },
        previousRow() {
            this.selection--;
            this.selection = updateSelection(this.selection, this.maxSelectionIndex || 0);
        },
        setSelection(selection) {
            this.selection = selection;
        }
    }
};

export const FormMixin = {
    props: {
        index: {default: false, required: true},
        selection: {default: 0, required: true},
        selectionIndex: {default: 0, required: true}
    },
    computed: {
        active() {
            return this.selectionIndex === this.selection;
        },
        showActive() {
            return this.active && !this.isTouch;
        }
    },
    methods: {
        focus() {
            this.$store.dispatch("selectRow", this.index);
            this.$emit("select", this.selectionIndex);
            this.focusCmdInput();
        }
    }
};
