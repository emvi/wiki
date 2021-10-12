<template>
    <emvi-list :title="$t('title')"
               icon="article"
               card="emvi-article-card"
               :perform-search="search"
               :disable-search="true"></emvi-list>
</template>

<script>
    import {ArticleService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList},
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    ArticleService.getDrafts(filter.offset, cancelToken)
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
            "title": "Drafts"
        },
        "de": {
            "title": "Entwürfe"
        }
    }
</i18n>
