<template>
    <div>
        <emvi-cmd-submenu icon="article" :enter="enter" :index="0" v-on:enter="createArticle" :disabled="isReadOnly">
            {{$t("submenu_article")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="list" view="list" :enter="enter" :index="1" :disabled="isReadOnly">
            {{$t("submenu_list")}}
        </emvi-cmd-submenu>
        <emvi-cmd-submenu icon="group" view="group" :enter="enter" :index="2" :disabled="!isExpert">
            {{$t("submenu_group")}}
            <span v-if="!isExpert">{{$t("requires_expert")}}</span>
        </emvi-cmd-submenu>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../util";
    import emviCmdOrganizationHeader from "../content/organization-header.vue";
    import emviCmdSubmenu from "../content/submenu.vue";

    export default {
        components: {emviCmdOrganizationHeader, emviCmdSubmenu},
        props: ["enter", "tab", "esc", "up", "down"],
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
            },
            tab(tab) {
                if(tab) {
                    if(!tab.shiftKey) {
                        this.$store.dispatch("selectNextRow");
                    }
                    else {
                        this.$store.dispatch("selectPreviousRow");
                    }
                }
            },
            esc(esc) {
                if(esc) {
                    this.$store.dispatch("popColumn");
                }
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
            }
        },
        methods: {
            createArticle() {
                this.$router.push("/edit");
                this.$store.dispatch("resetCmd");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "submenu_article": "Create a new Article",
            "submenu_list": "Create a new list",
            "submenu_group": "Create a new article",
            "requires_expert": "(requires an Expert organization)"
        },
        "de": {
            "submenu_article": "Neuen Artikel anlegen",
            "submenu_list": "Neue Liste anlegen",
            "submenu_group": "Neue Gruppe anlegen",
            "requires_expert": "(benötigt eine Expert Organisation)"
        }
    }
</i18n>
