<template>
    <emvi-list :title="$t('title')"
               icon="list"
               card="emvi-list-card"
               :perform-search="search"
               :disable-search="true">
        <emvi-private-tabs></emvi-private-tabs>
    </emvi-list>
</template>

<script>
    import {ArticlelistService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList, emviPrivateTabs} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList, emviPrivateTabs},
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    ArticlelistService.getPrivateLists(filter.offset, cancelToken)
                        .then(lists => {
                            resolve({results: lists, count: lists.length === 0 ? 0 : 9999999});
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
            "title": "Private"
        },
        "de": {
            "title": "Privat"
        }
    }
</i18n>
