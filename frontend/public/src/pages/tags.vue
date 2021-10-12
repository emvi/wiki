<template>
    <emvi-list :title="$t('title')"
               icon="tag"
               card="emvi-tag-card"
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
                    SearchService.findTags(query, filter, cancelToken)
                        .then(data => {
                            resolve({results: data.tags, count: data.count});
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
            "title": "Tags"
        },
        "de": {
            "title": "Tags"
        }
    }
</i18n>
