<template>
    <div class="input">
        <span>{{command}}</span>
        <input type="text"
               autocomplete="off"
               :placeholder="$t('cmd_placeholder')"
               ref="input"
               id="cmd-input"
               :value="command"
               v-on:input="e => command = e.target.value"
               v-on:keydown.enter.prevent.stop="enter"
               v-on:keydown.tab.prevent.stop="tab"
               v-on:keydown.esc.prevent.stop="esc"
               v-on:keydown.up.prevent.stop="up"
               v-on:keydown.down.prevent.stop="down"
               v-on:keydown.delete="del" />
        <input class="typeahead"
               v-model="typeahead"
               v-show="typeahead.startsWith(command)"
               tabindex="-1"
               readonly />    
    </div>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        data() {
            return {
                command: "",
                commandBefore: "",
                column: 0
            };
        },
        computed: {
            ...mapGetters(["cmd", "typeahead", "columns"])
        },
        watch: {
            cmd(cmd) {
                this.setCommand(cmd);
            },
            columns(columns) {
                if(columns === 0) {
                    this.command = "";
                }

                if(columns < this.column && (!document.activeElement || document.activeElement === document.body)) {
                    this.focus();
                }

                this.column = columns;
            },
            command() {
                this.update();
            }
        },
        mounted() {
            this.setCommand(this.cmd);
        },
        methods: {
            setCommand(cmd) {
                this.commandBefore = cmd;
                this.command = cmd;
                this.focus();
            },
            enter(e) {
                if(this.command.trim() !== "") {
                    this.$emit("enter", e);
                }
            },
            tab(e) {
                this.$emit("tab", e);
            },
            esc(e) {
                this.$emit("esc", e);
            },
            up(e) {
                this.$emit("up", e);
            },
            down(e) {
                this.$emit("down", e);
            },
            del(e) {
                this.$emit("del", e);
            },
            update() {
                let cmd = this.command.trim();

                if(cmd !== this.commandBefore) {
                    this.commandBefore = cmd;
                    this.$store.dispatch("resetCmd", cmd);
                }
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
            "cmd_placeholder": "Type a command or search..."
        },
        "de": {
            "cmd_placeholder": "Tippe einen Befehl oder suche..."
        }
    }
</i18n>
