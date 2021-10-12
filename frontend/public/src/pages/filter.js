export const FilterMixin = {
    data() {
        return {
            activeFilter: 0,
            maxFilterIndex: 0, // must be overwritten by component
            filterColumns: 1, // must be overwritten if page features a grid of filters
            enterEvent: false,
            tabEvent: false,
            escEvent: false,
            upEvent: false,
            downEvent: false,
            leftEvent: false,
            rightEvent: false
        };
    },
    methods: {
        enter(e) {
            this.enterEvent = e;
            this.$nextTick(() => {
                this.enterEvent = false;
            });
        },
        tab(e) {
            this.tabEvent = e;
            this.$nextTick(() => {
                this.tabEvent = false;
            });
        },
        esc(e) {
            this.escEvent = e;
            this.$nextTick(() => {
                this.escEvent = false;
            });
        },
        up(e) {
            if(e.shiftKey) {
                this.previousFilterRow();
            }
            else {
                this.upEvent = e;
                this.$nextTick(() => {
                    this.upEvent = false;
                });
            }
        },
        down(e) {
            if(e.shiftKey) {
                this.nextFilterRow();
            }
            else {
                this.downEvent = e;
                this.$nextTick(() => {
                    this.downEvent = false;
                });
            }
        },
        left(e) {
            if(e.shiftKey && this.filterColumns > 1) {
                this.previousFilter();
            }
            else {
                this.leftEvent = e;
                this.$nextTick(() => {
                    this.leftEvent = false;
                });
            }
        },
        right(e) {
            if(e.shiftKey && this.filterColumns > 1) {
                this.nextFilter();
            }
            else {
                this.rightEvent = e;
                this.$nextTick(() => {
                    this.rightEvent = false;
                });
            }
        },
        nextFilterRow() {
            this.activeFilter += this.filterColumns;

            if(this.activeFilter > this.maxFilterIndex) {
                this.activeFilter = 0;
            }
        },
        previousFilterRow() {
            this.activeFilter -= this.filterColumns;

            if(this.activeFilter < 0) {
                this.activeFilter = this.maxFilterIndex;
            }
        },
        nextFilter() {
            this.activeFilter++;

            if(this.activeFilter > this.maxFilterIndex) {
                this.activeFilter = 0;
            }
        },
        previousFilter() {
            this.activeFilter--;

            if(this.activeFilter < 0) {
                this.activeFilter = this.maxFilterIndex;
            }
        },
        toggleActiveFilter(index) {
            this.activeFilter = index;
        }
    }
};
