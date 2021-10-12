<template>
    <emvi-list :title="$t('title')"
               icon="group"
               card="emvi-group-card"
               :perform-search="search"
               :disable-search="true">
        <emvi-watch-tabs></emvi-watch-tabs>
    </emvi-list>
</template>

<script>
    import {ObserveService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList, emviWatchTabs} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviList, emviWatchTabs},
        methods: {
            search(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    ObserveService.getObserved(false, false, true, 0, 0, filter.offset, cancelToken)
                        .then(({groups}) => {
                            resolve({results: groups, count: groups.length === 0 ? 0 : 9999999});
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
            "title": "Watch List"
        },
        "de": {
            "title": "Beobachtungsliste"
        }
    }
</i18n>
