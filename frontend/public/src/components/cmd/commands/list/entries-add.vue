<template>
    <emvi-cmd-selection-search element="emvi-cmd-selection-article"
                               :perform-search="search"
                               :placeholder="$t('placeholder')"
                               :button="$t('label_add')"
                               icon="add"
                               :enter="enter"
                               :tab="tab"
                               :del="del"
                               :esc="esc"
                               :up="up"
                               :down="down"
                               v-on:save="save"></emvi-cmd-selection-search>
</template>

<script>
    import {addAttrToListElements} from "../../../../util";
    import {ArticlelistService, SearchService} from "../../../../service";
    import emviCmdSelectionSearch from "../../selection/search.vue";

    export default {
        components: {emviCmdSelectionSearch},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                filter.sort_title = "asc";

                return new Promise((resolve, reject) => {
                    SearchService.findArticles(query, filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results: addAttrToListElements(results, "type", "article"), count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            save(entities) {
                this.resetError();
                let listId = this.$store.state.page.meta.get("id");
                let ids = [];

                for(let i = 0; i < entities.length; i++) {
                    ids.push(entities[i].id);
                }

                ArticlelistService.addEntry(listId, ids)
                    .then(({results, count}) => {
                        this.$store.dispatch("success", this.$t("toast_added"));
                        this.$store.dispatch("popColumn");
                        this.$store.dispatch("setMeta", {key: "updateList", value: true});
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder": "Search articles...",
            "label_add": "Add Articles",
            "toast_added": "The articles have been added."
        },
        "de": {
            "placeholder": "Artikel suchen...",
            "label_add": "Artikel hinzufügen",
            "toast_added": "Die Artikel wurden hinzugefügt."
        }
    }
</i18n>
