<template>
    <div class="search">
        <div class="avatar size-40">
            <i class="icon icon-filter"></i>
        </div>
        <input type="text"
               autocomplete="off"
               ref="input"
               :placeholder="$t('placeholder')"
               :value="query"
               v-on:input="e => query = e.target.value"
               v-on:keydown.esc.stop.prevent="reset" />
    </div>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        data() {
            return {
                query: "",
                queryBefore: ""
            };
        },
        computed: {
            ...mapGetters(["selection", "cmdOpen"])
        },
        watch: {
            selection(index) {
                if(index === 0) {
                    this.focus();
                }
                else {
                    this.$refs.input.blur();
                }
            },
            cmdOpen(open) {
                if(this.selection === 0 && !open) {
                    this.focus();
                }
            },
            query() {
                this.update();
            }
        },
        mounted() {
            this.focus();
        },
        methods: {
            update() {
                if(this.query !== this.queryBefore) {
                    this.queryBefore = this.query;
                    this.$emit("update", this.query);
                }
            },
            reset() {
                if(this.queryBefore !== "") {
                    this.$emit("update", "");
                }

                this.query = "";
                this.queryBefore = "";
            },
            focus() {
                if(!this.isTouch) {
                    this.$refs.input.focus();
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder": "Type to filter results..."
        },
        "de": {
            "placeholder": "Tippe um die Ergebnisse zu filtern..."
        }
    }
</i18n>
