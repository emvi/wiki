<template>
    <emvi-cmd-selection element="emvi-cmd-selection-article-history"
                        :perform-search="search"
                        disable-search="true"
                        :enter="enter"
                        :tab="tab"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import emviCmdSelection from "../../selection/selection.vue";
    import {ArticleService} from "../../../../service";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        data() {
            return {
                loaded: false
            };
        },
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    let id = this.$store.state.page.meta.get("id");
                    let langId = this.$store.state.page.meta.get("langId");

                    ArticleService.getArticleHistory(id, langId, filter.offset, cancelToken)
                        .then(({history, count}) => {
                            if(history.length && !this.loaded) {
                                this.loaded = true;
                                history[0].is_latest = true;
                            }

                            // stop loading more after the latest three versions
                            if(!this.isExpert) {
                                count = 0;
                            }

                            resolve({results: history, count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>
