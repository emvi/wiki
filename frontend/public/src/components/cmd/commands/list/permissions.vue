<template>
    <emvi-cmd-selection element="emvi-cmd-selection-list-member"
                        :perform-search="search"
                        add-view="list-permissions-add"
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
    import {ArticlelistService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                let id = this.$store.state.page.meta.get("id");

                return new Promise((resolve, reject) => {
                    ArticlelistService.getMember(id, filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results, count});
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
            "add_text": "Add Member"
        },
        "de": {
            "add_text": "Mitglied hinzufügen"
        }
    }
</i18n>
