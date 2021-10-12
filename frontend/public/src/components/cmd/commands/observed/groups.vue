<template>
    <emvi-cmd-selection element="emvi-cmd-selection-observed-group"
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
                    ObserveService.getObserved(false, false, true, 0, 0, filter.offset, cancelToken)
                        .then(({groups}) => {
                            let results = [];

                            for(let i = 0; i < groups.length; i++) {
                                results.push(groups[i]);
                            }

                            resolve({results: addAttrToListElements(results, "type", "group"), count: groups.length === 0 ? 0 : 9999999});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>
