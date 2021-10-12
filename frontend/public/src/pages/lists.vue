<template>
    <emvi-list :title="$t('title')"
               icon="list"
               card="emvi-list-card"
               :perform-search="search"></emvi-list>
</template>

<script>
    import {SearchService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList},
        methods: {
            search(query, filter, cancelToken) {
                filter.sort_name = "asc";

                return new Promise((resolve, reject) => {
                    SearchService.findArticleLists(query, filter, cancelToken)
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
            "title": "Lists"
        },
        "de": {
            "title": "Listen"
        }
    }
</i18n>
