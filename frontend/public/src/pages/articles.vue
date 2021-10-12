<template>
    <emvi-list :title="$t('title')"
               icon="article"
               card="emvi-article-card"
               :perform-search="search"></emvi-list>
</template>

<script>
    import {SearchService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {
            emviList
        },
        methods: {
            search(query, filter, cancelToken) {
                filter.sort_title = "asc";

                return new Promise((resolve, reject) => {
                    SearchService.findArticles(query, filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results, count});
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
            "title": "Articles"
        },
        "de": {
            "title": "Artikel"
        }
    }
</i18n>
