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
                               v-on:save="add"></emvi-cmd-selection-search>
</template>

<script>
    import {SearchService} from "../../../../service";
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
            add(entities) {
                let access = [];

                for(let i = 0; i < entities.length; i++) {
                    if(entities[i].type === "user") {
                        access.push({user: entities[i], write: true});
                    }
                    else {
                        access.push({group: entities[i], write: true});
                    }
                }

                this.$store.dispatch("setMeta", {key: "setAccess", value: access});
                this.$store.dispatch("popColumn");
                this.$store.dispatch("success", this.$t("toast_added"));
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
