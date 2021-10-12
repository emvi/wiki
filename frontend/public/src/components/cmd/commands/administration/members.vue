<template>
    <emvi-cmd-selection element="emvi-cmd-selection-member"
                        :perform-search="search"
                        :placeholder="$t('placeholder')"
                        add-view="administration-members-add"
                        :add-text="$t('add_text')"
                        :enter="enter"
                        :tab="tab"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {addAttrToListElements} from "../../../../util";
    import {SearchService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                filter.sort_lastname = "asc";
                filter.sort_firstname = "asc";

                return new Promise((resolve, reject) => {
                    SearchService.findUser(query, filter, cancelToken)
                        .then(({user, count}) => {
                            resolve({results: addAttrToListElements(user, "type", "user"), count});
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
            "placeholder": "Filter Members...",
            "add_text": "Invite Members"
        },
        "de": {
            "placeholder": "Mitglieder filtern...",
            "add_text": "Mitglieder einladen"
        }
    }
</i18n>
