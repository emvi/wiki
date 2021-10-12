<template>
    <emvi-cmd-selection element="emvi-cmd-selection-observed-article"
                        :disable-search="true"
                        :perform-search="search"
                        :enter="enter"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {ObserveService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";
    import {addAttrToListElements} from "../../../../util";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    ObserveService.getObserved(true, false, false, filter.offset, 0, 0, cancelToken)
                        .then(({articles}) => {
                            let results = [];

                            for(let i = 0; i < articles.length; i++) {
                                results.push(articles[i]);
                            }

                            resolve({results: addAttrToListElements(results, "type", "article"), count: articles.length === 0 ? 0 : 9999999});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>
