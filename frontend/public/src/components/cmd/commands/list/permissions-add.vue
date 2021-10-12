<template>
    <emvi-cmd-selection-search element="emvi-cmd-selection-user-group"
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
    import {ArticlelistService, SearchService} from "../../../../service";
    import emviCmdSelectionSearch from "../../selection/search.vue";

    export default {
        components: {emviCmdSelectionSearch},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    SearchService.findUserAndUserGroup(query, cancelToken)
                        .then(results => {
                            resolve({results, count: 0});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            save(entities) {
                this.resetError();
                let listId = this.$store.state.page.meta.get("id");
                let user = [];
                let groups = [];

                for(let i = 0; i < entities.length; i++) {
                    if(entities[i].type === "user") {
                        user.push(entities[i].id);
                    }
                    else {
                        groups.push(entities[i].id);
                    }
                }

                ArticlelistService.addMember(listId, user, groups)
                    .then(({results, count}) => {
                        this.$store.dispatch("success", this.$t("toast_added"));
                        this.$store.dispatch("popColumn");
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
            "placeholder": "Search members or groups...",
            "label_add": "Add Members",
            "toast_added": "The members have been added."
        },
        "de": {
            "placeholder": "Mitglieder oder Gruppen suchen...",
            "label_add": "Mitglieder hinzufügen",
            "toast_added": "Die Mitglieder wurden hinzugefügt."
        }
    }
</i18n>
