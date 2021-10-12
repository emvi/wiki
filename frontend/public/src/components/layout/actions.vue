<template>
    <div :class="{'drawer right': true, 'expand': !hidden, 'no-transition': noTransition}"
         v-show="commands.length"
         v-if="member.show_action_buttons"
         ref="actions">
        <i class="toggle icon icon-flash cursor-pointer" v-on:click="toggle"></i>
        <div class="drawer-content focus-ring">
            <button v-for="cmd in commands"
                    :key="cmd.id"
                    :class="{disabled: cmd.isDisabled}"
                    :title="cmd.expert && !isExpert ? $t('expert'): ''"
                    v-on:click="run(cmd)">
                <i :class="getIconClass(cmd.icon)"></i>
                {{getLabel(cmd)}}
            </button>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {capitalize} from "../../util";
    import {filterCommands} from "../cmd/commands";

    export default {
        props: {
            hide: {default: false}
        },
        data() {
            return {
                mouseHandler: null,
                commands: [],
                noTransition: true,
                hidden: true
            };
        },
        computed: {
            ...mapGetters(["member", "metaUpdate"]),
            page() {
                return this.$route.name;
            }
        },
        watch: {
            page() {
                this.filterCommands();
            },
            metaUpdate() {
                this.filterCommands();
            }
        },
        mounted() {
            this.hidden = this.hide || this.$mq <= 1440;

            if(this.hidden) {
                this.bindMouse();
            }

            this.filterCommands();
        },
        beforeDestroy() {
            if(this.mouseHandler) {
                this.unbindMouse();
            }
        },
        methods: {
            bindMouse() {
                this.mouseHandler = document.addEventListener("mouseup", e => {
                    if(this.$refs.actions && e.target !== this.$refs.actions && !this.$refs.actions.contains(e.target)) {
                        this.hidden = true;
                    }
                });
            },
            unbindMouse() {
                document.removeEventListener("mouseup", this.mouseHandler);
            },
            capitalize(name) {
                return capitalize(name);
            },
            getIconClass(icon) {
                let c = {icon: true};
                c[`icon-${icon} size-40`] = true;
                return c;
            },
            getLabel(cmd) {
                return cmd.button_label ? cmd.button_label[this.locale] : capitalize(cmd.command[this.locale][0]);
            },
            filterCommands() {
                let filtered = filterCommands(this.page, "", this.locale);
                let commands = [];

                // filter for commands that are specific to this page
                for(let i = 0; i < filtered.length; i++) {
                    if(filtered[i].pageNames || filtered[i].buttonPageNames && filtered[i].buttonPageNames.indexOf(this.page) > -1) {
                        filtered[i].isDisabled = filtered[i].disabled && filtered[i].disabled(this);
                        commands.push(filtered[i]);
                    }
                }

                this.commands = commands;
            },
            run(cmd) {
                if(cmd.expert && !this.isExpert) {
                    this.$router.push("/billing");
                    return;
                }

                if(cmd.isDisabled) {
                    return;
                }

                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", "/"+cmd.command[this.locale][0]);
                    cmd.run(this);
                });
            },
            toggle() {
                this.noTransition = false;
                this.hidden = !this.hidden;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "expert": "Requires an Expert organization."
        },
        "de": {
            "expert": "Benötigt eine Expert Organisation."
        }
    }
</i18n>
