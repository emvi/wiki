<template>
    <emvi-cmd-selection element="emvi-cmd-selection-article-member"
                        :perform-search="search"
                        add-view="article-permissions-member-add"
                        :add-text="$t('add_text')"
                        disable-search="true"
                        :update-results="updateResults"
                        :enter="enter"
                        :tab="tab"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {mapGetters} from "vuex";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        data() {
            return {
                updateResults: 0,
                count: 0
            };
        },
        computed: {
            ...mapGetters(["metaUpdate"])
        },
        watch: {
            metaUpdate() {
                if(this.$store.state.page.meta.get("access").length !== this.count) {
                    this.updateResults++;
                }
            }
        },
        methods: {
            search() {
                return new Promise((resolve, reject) => {
                    let member = this.$store.state.page.meta.get("access");
                    this.count = member.length;
                    let results = [];

                    for(let i = 0; i < member.length; i++) {
                        if(member[i].user) {
                            member[i].user.type = "user";
                            member[i].user.write = member[i].write;
                            results.push(member[i].user);
                        }
                        else {
                            member[i].group.type = "group";
                            member[i].group.write = member[i].write;
                            results.push(member[i].group);
                        }
                    }

                    resolve({results, count: 0});
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "add_text": "Add Members"
        },
        "de": {
            "add_text": "Mitglieder hinzufügen"
        }
    }
</i18n>
