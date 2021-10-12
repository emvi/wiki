<template>
    <emvi-list :title="$t('title')"
               icon="article"
               card="emvi-article-card"
               :perform-search="search"
               :disable-search="true">
        <emvi-private-tabs></emvi-private-tabs>
    </emvi-list>
</template>

<script>
    import {ArticleService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList, emviPrivateTabs} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList, emviPrivateTabs},
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    ArticleService.getPrivateArticles(filter.offset, cancelToken)
                        .then(articles => {
                            resolve({results: articles, count: articles.length === 0 ? 0 : 9999999});
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
