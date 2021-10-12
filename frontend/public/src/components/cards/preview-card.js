import {scrollToTop, scrollToTopArea} from "../../util";

export const CardMixin = {
    props: {
        entity: {default: () => {return {}}},
        index: {default: 0},
        active: {default: false},
        scrollArea: {default: ""},
        enter: {default: false},
        tab: {default: false},
        up: {default: false},
        down: {default: false},
        query: {default: ""}
    },
    data() {
        return {
            showPreview: false
        };
    },
    watch: {
        active(active) {
            if(!active && this.showPreview) {
                this.showPreview = false;
            }
        }
    },
    methods: {
        scroll() {
            this.$nextTick(() => {
                if(this.scrollArea) {
                    scrollToTopArea(this.$refs.card.$refs.card.$refs.card, this.scrollArea);
                }
                else {
                    scrollToTop(this.$refs.card.$refs.card.$refs.card);
                }
            });
        }
    }
};
