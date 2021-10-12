<template>
    <emvi-cmd-selection-search element="emvi-cmd-selection-user-group"
                               :preselected="selected"
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
    import {mapGetters} from "vuex";
    import {SearchService} from "../../../service";
    import emviCmdSelectionSearch from "../selection/search.vue";

    export default {
        components: {emviCmdSelectionSearch},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        data() {
            return {
                selected: []
            };
        },
        computed: {
            ...mapGetters(["cmdMeta"])
        },
        beforeMount() {
            this.init();
        },
        methods: {
            init() {
                if(this.cmdMeta.get("members")) {
                    this.selected = this.cmdMeta.get("members");
                }
            },
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
                this.$store.dispatch("popColumn");
                this.$store.dispatch("setCmdMeta", {key: "members", value: entities});
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder": "Search members or groups...",
            "label_add": "Select Members"
        },
        "de": {
            "placeholder": "Mitglieder oder Gruppen suchen...",
            "label_add": "Mitglieder auswählen"
        }
    }
</i18n>
