<template>
    <transition name="slide-up-far">
        <div v-if="toastMessage" v-bind:class="{'toast--blue': toastType === 'blue', 'toast--green': toastType === 'green', 'toast--orange': toastType === 'orange', 'toast--red': toastType === 'red'}">
            <i class="icon icon-error toast--icon"></i>
            <div class="toast--text">
                <span v-if="toastException">An unexpected error occurred ({{new Date() | moment("YYYY-MM-DD HH:mm:ss")}}):</span>
                {{toastMessage}}
            </div>
            <i class="icon icon-close-small toast--close" v-on:click="close"></i>
        </div>
    </transition>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        computed: {
            ...mapGetters(["toastException", "toastMessage", "toastType"])
        },
        methods: {
            close() {
                this.$store.dispatch("resetToast");
            }
        }
    }
</script>
