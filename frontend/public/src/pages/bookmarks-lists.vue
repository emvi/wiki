<template>
    <emvi-list :title="$t('title')"
               icon="list"
               card="emvi-list-card"
               :perform-search="search"
               :disable-result-count="true"
               :disable-search="true">
        <emvi-bookmark-tabs></emvi-bookmark-tabs>
    </emvi-list>
</template>

<script>
    import {BookmarkService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList, emviBookmarkTabs} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList, emviBookmarkTabs},
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    BookmarkService.getBookmarks(false, true, 0, filter.offset, cancelToken)
                        .then(({lists}) => {
                            let results = [];

                            for(let i = 0; i < lists.length; i++) {
                                results.push(lists[i].article_list);
                            }

                            resolve({results, count: lists.length === 0 ? 0 : 9999999});
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
            "title": "Bookmarks"
        },
        "de": {
            "title": "Lesezeichen"
        }
    }
</i18n>
