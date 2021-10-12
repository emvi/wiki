<template>
    <emvi-cmd-selection element="emvi-cmd-selection-list-article"
                        :perform-search="search"
                        add-view="list-entries-add"
                        :add-text="$t('add_text')"
                        disable-search="true"
                        :enter="enter"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {addAttrToListElements} from "../../../../util";
    import {ArticlelistService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                let id = this.$store.state.page.meta.get("id");

                return new Promise((resolve, reject) => {
                    ArticlelistService.getEntries(id, filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results: addAttrToListElements(results, "type", "article"), count});
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
            "add_text": "Add Articles"
        },
        "de": {
            "add_text": "Artikel hinzufügen"
        }
    }
</i18n>
