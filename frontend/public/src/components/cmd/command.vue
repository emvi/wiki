<template>
    <div>
        <emvi-cmd-result v-for="(command, index) in commands"
                         :key="index"
                         :icon="command.icon"
                         :expert="command.expert"
                         :requires-admin="command.requires_admin"
                         :requires-mod="command.requires_mod"
                         :disable-func="command.disabled"
                         :index="index"
                         v-on:click="exec(index)">
            <p>
                <template v-if="!command.command[locale][0].startsWith(query)">/{{findCommandAlias(command)}} -&gt; </template>
                /{{command.command[locale][0]}}
            </p>
            <small>
                {{command.description[locale]}}
            </small>
        </emvi-cmd-result>
        <emvi-cmd-empty v-show="!commands.length"></emvi-cmd-empty>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {filterCommands, findCommandAlias} from "./commands";
    import {getQueryFromCmd, updateSelectedRow} from "./util";
    import emviCmdResult from "./content/result.vue";
    import emviCmdEmpty from "./content/empty.vue";

    export default {
        components: {emviCmdResult, emviCmdEmpty},
        props: ["up", "down", "enter", "tab", "esc"],
        data() {
            return {
                commands: []
            };
        },
        computed: {
            ...mapGetters(["cmd", "row"]),
            query() {
                return getQueryFromCmd(this.cmd);
            }
        },
        watch: {
            cmd(cmd) {
                this.filterCommands(cmd);
            },
            row(row) {
                updateSelectedRow(row, this.commands.length, this.$store);
                this.$store.dispatch("typeahead", this.typeahead());
            },
            up(up) {
                if(up) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down) {
                    this.$store.dispatch("selectNextRow");
                }
            },
            enter(enter) {
                if(enter) {
                    this.exec();
                }
            },
            tab(tab) {
                if(tab) {
                    let cmd = this.autocomplete();

                    if(cmd) {
                        this.$store.dispatch("setCmd", cmd);
                    }
                }
            },
            esc(esc) {
                if(esc) {
                    this.$store.dispatch("popColumn");
                }
            }
        },
        mounted() {
            this.filterCommands(this.cmd);
        },
        methods: {
            filterCommands(cmd) {
                this.commands = filterCommands(this.$route.name, cmd, this.locale);
                this.$store.dispatch("typeahead", this.typeahead());
            },
            typeahead() {
                if(!this.commands.length) {
                    return "";
                }

                let cmd = this.commands[this.row];
                return `/${findCommandAlias(cmd, this.cmd, this.locale)}`;
            },
            autocomplete() {
                if(!this.commands.length) {
                    return "";
                }

                let cmd = this.commands[this.row];
                this.$store.dispatch("selectRow", 0);
                return `/${findCommandAlias(cmd, this.cmd, this.locale)}`;
            },
            exec(index) {
                if(!this.commands.length) {
                    return;
                }

                if(index === undefined) {
                    index = this.row;
                }

                let cmd = this.commands[index];

                if((cmd.expert && !this.isExpert) ||
                    (cmd.requires_admin && !this.isAdmin) ||
                    (cmd.requires_mod && !this.isAdmin && !this.isMod) ||
                    (cmd.disabled && cmd.disabled(this))) {
                    return;
                }

                cmd.run(this);
            },
            findCommandAlias(cmd) {
                return findCommandAlias(cmd, this.cmd, this.locale);
            }
        }
    }
</script>
