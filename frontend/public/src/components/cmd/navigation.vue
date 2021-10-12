<template>
    <div>
        <emvi-cmd-result v-for="(nav, index) in navigation"
                         :key="nav.page[locale][0]"
                         :icon="nav.icon"
                         :index="index"
                         v-on:click="navigate(index)">
            <p>
                <template v-if="!nav.page[locale][0].startsWith(query)">.{{findNavigationAlias(nav)}} -&gt; </template>
                .{{nav.page[locale][0]}}
            </p>
            <small>
                {{nav.title[locale]}}
            </small>
        </emvi-cmd-result>
        <emvi-cmd-empty v-show="!navigation.length"></emvi-cmd-empty>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {filterNavigation, findNavigationAlias} from "./navigation";
    import {getQueryFromCmd, updateSelectedRow} from "./util";
    import emviCmdResult from "./content/result.vue";
    import emviCmdEmpty from "./content/empty.vue";

    export default {
        components: {emviCmdResult, emviCmdEmpty},
        props: ["up", "down", "enter", "tab", "esc"],
        data() {
            return {
                navigation: []
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
                this.filterNavigation(cmd);
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
            row(row) {
                updateSelectedRow(row, this.navigation.length, this.$store);
                this.$store.dispatch("typeahead", this.typeahead());
            },
            enter(enter) {
                if(enter) {
                    this.navigate();
                    this.$store.dispatch("resetCmd");
                    this.$store.dispatch("select", 0); // reset page selection
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
            this.filterNavigation(this.cmd);
        },
        methods: {
            filterNavigation(cmd) {
                this.navigation = filterNavigation(this.$route.name, cmd, this.locale, this.isAdmin);
                this.$store.dispatch("typeahead", this.typeahead());
            },
            typeahead() {
                if(!this.navigation.length) {
                    return "";
                }

                let nav = this.navigation[this.row];
                return `.${findNavigationAlias(nav, this.cmd, this.locale)}`;
            },
            autocomplete() {
                if(!this.navigation.length) {
                    return "";
                }

                let nav = this.navigation[this.row];
                return `.${findNavigationAlias(nav, this.cmd, this.locale)}`;
            },
            navigate(index) {
                if(!this.navigation.length) {
                    return;
                }

                if(index === undefined) {
                    index = this.row;
                }

                this.$router.push(this.navigation[index].path);
                this.$store.dispatch("resetCmd");
            },
            findNavigationAlias(nav) {
                return findNavigationAlias(nav, this.cmd, this.locale);
            }
        }
    }
</script>
