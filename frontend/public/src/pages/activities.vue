<template>
    <emvi-list :title="$t('title')"
               icon="notification"
               card="emvi-activity-card"
               :perform-search="search"
               :disable-search="true"
               disable-result-count="true"></emvi-list>
</template>

<script>
    import {FeedService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList},
        methods: {
            search(query, filter, cancelToken) {
                filter.limit = 20;

                return new Promise((resolve, reject) => {
                    FeedService.getFeed(filter, cancelToken)
                        .then(results => {
                            resolve({results, count: results.length === 0 ? 0 : 9999999});
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
            "title": "Activities"
        },
        "de": {
            "title": "Aktivitäten"
        }
    }
</i18n>
