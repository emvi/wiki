<template>
    <div :class="toastClass" v-show="showToast" v-on:click="close">
        <i :class="iconClass"></i>
        <div class="message">{{toastMessage}}</div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        computed: {
            ...mapGetters(["showToast", "toastColor", "toastMessage"]),
            toastClass() {
                let c = {};
                c["toast no-select"] = true;
                c["green"] = this.toastColor === "green";
                c["red"] = this.toastColor === "red";
                return c;
            },
            iconClass() {
                let c = {};
                c["icon"] = true;
                c["icon-check green-100"] = this.toastColor === "green";
                c["icon-error red-100"] = this.toastColor === "red";
                return c;
            }
        },
        methods: {
            close() {
                this.$store.dispatch("closeToast");
            }
        }
    }
</script>
