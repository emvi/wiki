<template>
    <emvi-cmd-selection element="emvi-cmd-selection-group-member"
                        :perform-search="search"
                        add-view="group-member-add"
                        :add-text="$t('add_text')"
                        disable-search="true"
                        :enter="enter"
                        :tab="tab"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {UsergroupService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                let id = this.$store.state.page.meta.get("id");

                return new Promise((resolve, reject) => {
                    UsergroupService.getMember(id, filter, cancelToken)
                        .then(({results, count}) => {
                            let members = [];

                            for(let i = 0; i < results.length; i++) {
                                results[i].user.type = "user";
                                results[i].user.member_id = results[i].id;
                                results[i].user.is_moderator = results[i].is_moderator;
                                members.push(results[i].user);
                            }

                            resolve({results: members, count});
                        })
                        .catch(e => {
                            reject(e);
                        });
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
