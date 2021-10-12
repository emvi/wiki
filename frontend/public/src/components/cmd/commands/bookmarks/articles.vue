<template>
    <emvi-cmd-selection element="emvi-cmd-selection-bookmarked-article"
                        :disable-search="true"
                        :perform-search="search"
                        :enter="enter"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {BookmarkService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "del", "esc", "up", "down"],
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    BookmarkService.getBookmarks(true, false, filter.offset, 0, cancelToken)
                        .then(({articles}) => {
                            let results = [];

                            for(let i = 0; i < articles.length; i++) {
                                articles[i].article.type = "article";
                                results.push(articles[i].article);
                            }

                            resolve({results: results, count: articles.length === 0 ? 0 : 9999999});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>
