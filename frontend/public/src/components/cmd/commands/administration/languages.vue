<template>
    <emvi-cmd-selection element="emvi-cmd-selection-language"
                        :perform-search="search"
                        :placeholder="$t('placeholder')"
                        add-view="administration-languages-add"
                        :add-text="$t('add_text')"
                        :enter="enter"
                        :tab="tab"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {LangService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "esc", "up", "down"],
        methods: {
            search(query) {
                query = query.trim().toLowerCase();

                return new Promise((resolve, reject) => {
                    this.getLangs()
                        .then(langs => {
                            if(query === "") {
                                resolve({results: langs, count: 0});
                            }
                            else {
                                let filtered = [];

                                for(let i = 0; i < langs.length; i++) {
                                    if(langs[i].code.includes(query) || langs[i].name.toLowerCase().includes(query)) {
                                        filtered.push(langs[i]);
                                    }
                                }

                                resolve({results: filtered, count: 0});
                            }
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            getLangs() {
                return new Promise((resolve, reject) => {
                    LangService.getLangs()
                        .then(langs => {
                            resolve(langs);
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
            "placeholder": "Filter languages...",
            "add_text": "Add Language"
        },
        "de": {
            "placeholder": "Sprachen filtern...",
            "add_text": "Sprache hinzufügen"
        }
    }
</i18n>
