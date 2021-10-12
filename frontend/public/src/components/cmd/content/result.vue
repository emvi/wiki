<template>
    <div :class="resultClass" v-on:click.stop.prevent="click" ref="result">
        <emvi-cmd-icon :icon="icon"></emvi-cmd-icon>
        <span class="item">
            <slot></slot>
        </span>
        <span class="label" v-if="permissionAdmin">
            {{$t("administrator")}}
        </span>
        <span class="label" v-if="permissionMod">
            {{$t("moderator")}}
        </span>
        <span class="link cursor-pointer" v-if="requiresExpert" v-on:click.stop.prevent="toBillingPage">
            {{$t("upgrade")}}
        </span>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {scrollIntoViewArea} from "../../../util";
    import emviCmdIcon from "../content/icon.vue";

    export default {
        components: {emviCmdIcon},
        props: {
            icon: "",
            expert: {default: false},
            requiresAdmin: {default: false},
            requiresMod: {default: false},
            disableReadOnly: {default: false},
            disableFunc: null,
            index: {default: 0}
        },
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.index === this.row;
            },
            resultClass() {
                return {
                    "entry cursor-pointer": true,
                    active: this.active && !this.isTouch,
                    disabled: this.disabled,
                    double: this.isMobile
                };
            },
            permissionAdmin() {
                return this.requiresAdmin && !this.isAdmin;
            },
            permissionMod() {
                return this.requiresMod && !this.isAdmin && !this.isMod;
            },
            requiresExpert() {
                return this.expert && !this.isExpert;
            },
            disabled() {
                return this.permissionAdmin ||
                    this.permissionMod ||
                    this.requiresExpert ||
                    this.disableFunc && this.disableFunc(this);
            }
        },
        watch: {
            row() {
                if(this.active) {
                    scrollIntoViewArea(this.$refs.result, document.getElementById("cmd-results"));
                }
            }
        },
        methods: {
            click() {
                this.$emit('click');
                this.focusCmdInput();
            },
            toBillingPage() {
                if(this.$route.path !== "/billing") {
                    this.$router.push("/billing");
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "administrator": "Administrator",
            "moderator": "Moderator",
            "upgrade": "Upgrade to Expert"
        },
        "de": {
            "administrator": "Administrator",
            "moderator": "Moderator",
            "upgrade": "Upgrade auf Expert"
        }
    }
</i18n>
