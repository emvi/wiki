<template>
    <emvi-cmd-selection element="emvi-cmd-selection-article-translation"
                        :perform-search="search"
                        disable-search="true"
                        :enter="enter"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import emviCmdSelection from "../../selection/selection.vue";
    import {LangService} from "../../../../service";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "esc", "up", "down"],
        data() {
            return {
                loaded: false
            };
        },
        methods: {
            search() {
                return new Promise((resolve, reject) => {
                    LangService.getLangs()
                        .then(results => {
                            resolve({results, count: results.length});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>
